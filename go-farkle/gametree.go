package farkle

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"iter"
	"math"
	"os"
	"runtime"
	"sync"

	"github.com/bsm/extsort"
	"github.com/golang/glog"
)

// Action is the choice made by a player after rolling.
// A zero Action is a Farkle.
type Action struct {
	HeldDiceID      uint16
	ContinueRolling bool
}

func (a Action) String() string {
	if a.HeldDiceID == 0 && !a.ContinueRolling {
		return "FARKLE!"
	}

	roll := rollsByID[a.HeldDiceID]
	contStr := "Stop"
	if a.ContinueRolling {
		contStr = "Continue"
	}
	return fmt.Sprintf("{Held: %s, %s}", roll, contStr)
}

func ApplyAction(state GameState, action Action) GameState {
	trickScore := scoreCache[action.HeldDiceID]
	newScore := state.ScoreThisRound + trickScore
	if newScore < state.ScoreThisRound {
		newScore = math.MaxUint8 // Overflow
	}
	state.ScoreThisRound = newScore
	if trickScore == 0 { // Farkle
		state.ScoreThisRound = 0
	}

	numDiceHeld := rollNumDice[action.HeldDiceID]
	if numDiceHeld > state.NumDiceToRoll {
		panic(fmt.Errorf("illegal action %+v applied to state %+v: "+
			"held %d dice but only had %d to roll",
			action, state, numDiceHeld, state.NumDiceToRoll))
	}
	state.NumDiceToRoll -= numDiceHeld
	if state.NumDiceToRoll == 0 {
		state.NumDiceToRoll = MaxNumDice
	}

	if !action.ContinueRolling {
		currentScore := state.PlayerScores[0]
		newScore := currentScore + state.ScoreThisRound
		if newScore < currentScore {
			newScore = math.MaxUint8 // Overflow
		}
		// Advance to next player by rotating the scores.
		copy(state.PlayerScores[:state.NumPlayers], state.PlayerScores[1:state.NumPlayers])
		state.PlayerScores[state.NumPlayers-1] = newScore
		state.ScoreThisRound = 0
		state.NumDiceToRoll = MaxNumDice
	}

	return state
}

// Find the action that maximizes current player win probability.
func SelectAction(state GameState, rollID uint16, db DB) (Action, [maxNumPlayers]float64) {
	var bestWinProb [maxNumPlayers]float64
	var bestAction Action
	notYetOnBoard := (state.PlayerScores[0] == 0)
	potentialActions := rollIDToPotentialActions[rollID]
	for _, action := range potentialActions {
		if state.ScoreThisRound == math.MaxUint8 && action.ContinueRolling {
			// Overflowed score this round. Our assumption is that this is unlikely.
			// Approximate the solution using the probability as if they stopped.
			action.ContinueRolling = false
		}

		newState := ApplyAction(state, action)
		if notYetOnBoard && !action.ContinueRolling && newState.PlayerScores[state.NumPlayers-1] < 500/incr {
			// Not a valid state: You must get at least 500 to get on the board.
			continue
		}

		pSubtree := db.Get(newState)
		if !action.ContinueRolling {
			// Probabilities are rotated since we advanced to the
			// next player in next state.
			pSubtree = unrotate(pSubtree, state.NumPlayers)
		}
		if pSubtree[0] > bestWinProb[0] {
			bestWinProb = pSubtree
			bestAction = action
		}
	}

	if len(potentialActions) == 0 {
		newState := ApplyAction(state, bestAction)
		pSubtree := db.Get(newState)
		bestWinProb = unrotate(pSubtree, state.NumPlayers)
	}

	return bestAction, bestWinProb
}

func unrotate(pWin [maxNumPlayers]float64, numPlayers uint8) [maxNumPlayers]float64 {
	var result [maxNumPlayers]float64
	copy(result[1:numPlayers], pWin[:numPlayers])
	result[0] = pWin[numPlayers-1]
	return result
}

var rollIDToPotentialActions = func() [][]Action {
	result := make([][]Action, len(rollIDToPotentialHolds))
	for rollID, holds := range rollIDToPotentialHolds {
		actions := make([]Action, 0, 2*len(holds))
		for _, holdOption := range holds {
			for _, continueRolling := range []bool{true, false} {
				actions = append(actions, Action{
					HeldDiceID:      rollToID[holdOption],
					ContinueRolling: continueRolling,
				})
			}
		}

		result[rollID] = actions
	}

	return result
}()

// Recalculate the value of all states in the given iterator,
// updating the value of each state in the database.
func UpdateAll(db DB, states iter.Seq2[uint16, GameState]) {
	// Recalculate all other states.
	var mx sync.RWMutex
	var wg sync.WaitGroup
	var workCh chan GameState
	numWorkers := runtime.NumCPU()
	currentDepth := uint16(0)
	for depth, state := range states {
		if depth != currentDepth {
			// Wait for previous depth to complete.
			close(workCh)
			wg.Wait()

			// Start up workers for next depth.
			glog.Infof("Processing game states with depth=%d", depth)
			workCh = make(chan GameState, numWorkers)
			wg.Add(numWorkers)
			for i := 0; i < numWorkers; i++ {
				go func() {
					updateWorker(db, workCh, &mx)
					wg.Done()
				}()
			}
		}

		workCh <- state
	}

	close(workCh)
	wg.Wait()
}

func updateWorker(db DB, workCh <-chan GameState, mx *sync.RWMutex) {
	// We batch updates to the database to reduce lock contention.
	batchSize := 1024 // Arbitrary, tunable
	batchStates := make([]GameState, batchSize)
	batchUpdates := make([][maxNumPlayers]float64, batchSize)
	for state := range workCh {
		var pWin [maxNumPlayers]float64
		if state.IsGameOver() {
			pWin = calcEndGameValue(state)
		} else {
			mx.RLock()
			pWin = calcStateValue(state, db)
			mx.RUnlock()
		}

		batchStates = append(batchStates, state)
		batchUpdates = append(batchUpdates, pWin)
		if len(batchStates) == cap(batchStates) {
			mx.Lock()
			for i, state := range batchStates {
				db.Put(state, batchUpdates[i])
			}
			mx.Unlock()
			batchStates = batchStates[:0]
			batchUpdates = batchUpdates[:0]
		}
	}

	mx.Lock()
	defer mx.Unlock()
	for i, state := range batchStates {
		db.Put(state, batchUpdates[i])
	}
}

func calcEndGameValue(state GameState) [maxNumPlayers]float64 {
	winningScore := state.HighestScore()
	winners := make([]int, 0, maxNumPlayers)
	for player, score := range state.PlayerScores[:state.NumPlayers] {
		if score == winningScore {
			winners = append(winners, player)
		}
	}

	// Not clear how ties should be considered in terms of "win probability".
	// We split the win amongst all players with the same score.
	p := 1.0 / float64(len(winners))
	var result [maxNumPlayers]float64
	for _, winner := range winners {
		result[winner] = p
	}

	return result
}

func calcStateValue(state GameState, db DB) [maxNumPlayers]float64 {
	var pWin [maxNumPlayers]float64
	for _, wRoll := range allRolls[state.NumDiceToRoll] {
		_, pSubgame := SelectAction(state, wRoll.ID, db)
		for i, p := range pSubgame[:state.NumPlayers] {
			pWin[i] += wRoll.Prob * p
		}
	}

	return pWin
}

// Save all game states from the given iterator to a file.
func SaveGameStates(states iter.Seq2[uint16, GameState], path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriterSize(f, 4*1024*1024)

	glog.Infof("Saving game states to: %s", path)
	buf := make([]byte, maxSizeOfGameState+2)
	i := 0
	for depth, state := range states {
		binary.LittleEndian.PutUint16(buf[:2], depth)
		n := state.SerializeTo(buf[1:])
		if _, err := w.Write(buf[:n+1]); err != nil {
			return err
		}

		i++
		if i%100000 == 0 {
			glog.Infof("...%d", i)
		}
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return f.Close()
}

// Return an iterator over all game states in the given file.
func IterGameStates(numPlayers int, path string) (iter.Seq2[uint16, GameState], error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return func(yield func(uint16, GameState) bool) {
		defer f.Close()
		r := bufio.NewReaderSize(f, 4*1024*1024)

		buf := make([]byte, numPlayers+3+2)
		for {
			_, err := io.ReadFull(r, buf)
			if err == io.EOF {
				break
			} else if err != nil {
				panic(fmt.Errorf("error reading game states: %w", err))
			}

			depth := binary.LittleEndian.Uint16(buf[:2])
			state := GameStateFromBytes(buf[2:])
			if !yield(depth, state) {
				break
			}
		}
	}, nil
}

// Return an iterator over all distinct game states and their depth in the game tree.
// Game states are sorted by depth in descending order such that end game states
// are enumerated before early game states.
func SortedGameStates(numPlayers int, workDir string) iter.Seq2[uint16, GameState] {
	sorter := extsort.New(&extsort.Options{
		WorkDir:    workDir,
		Compare:    compareGameStateDepth,
		BufferSize: 16 * 1024 * 1024, // 16 GiB
	})

	glog.Infof("Enumerating all %d %d-player game states",
		calcNumDistinctStates(numPlayers), numPlayers)
	i := 0
	for depth, gs := range allGameStates(numPlayers) {
		if depth > math.MaxUint16 {
			panic(fmt.Errorf("game state has depth %d > max uint8", depth))
		}

		key := make([]byte, 2)
		binary.LittleEndian.PutUint16(key, uint16(depth))
		value := make([]byte, numPlayers+3)
		n := gs.SerializeTo(value)
		sorter.Put(key, value[:n])

		i++
		if i%100000 == 0 {
			glog.Infof("...%d", i)
		}
	}

	glog.Info("Sorting game states by depth")
	iter, err := sorter.Sort()
	if err != nil {
		panic(fmt.Errorf("error sorting game states: %w", err))
	}

	return func(yield func(uint16, GameState) bool) {
		for iter.Next() {
			depth := binary.LittleEndian.Uint16(iter.Key())
			state := GameStateFromBytes(iter.Value())
			if !yield(depth, state) {
				break
			}
		}

		if err := iter.Err(); err != nil {
			panic(fmt.Errorf("error sorting game states: %w", err))
		}

		if err := iter.Close(); err != nil {
			panic(fmt.Errorf("error sorting game states: %w", err))
		}
	}
}

// Sort game states deeper in the tree before earlier states.
// i.e. end game -> initial state
func compareGameStateDepth(d1, d2 []byte) int {
	m := binary.LittleEndian.Uint16(d1)
	n := binary.LittleEndian.Uint16(d2)

	if m > n {
		return -1
	} else if m == n {
		return 0
	}

	return 1
}

// Return an iterator over all distinct game states, and their minimum
// depth in the game tree.
func allGameStates(numPlayers int) iter.Seq2[int, GameState] {
	return func(yield func(int, GameState) bool) {
		initialState := NewGameState(numPlayers)
		mask := newBitMask(calcNumDistinctStates(numPlayers))
		recursiveEnumerateStates(initialState, mask, 0, yield)
	}
}

func recursiveEnumerateStates(state GameState, mask *bitMask, depth int, yield func(int, GameState) bool) bool {
	gsID := state.ID()
	if mask.IsSet(gsID) {
		return true
	}

	mask.Set(gsID)
	if state.IsGameOver() {
		return yield(depth, state)
	}

	notYetOnBoard := (state.PlayerScores[0] == 0)
	for _, wRoll := range allRolls[state.NumDiceToRoll] {
		potentialActions := rollIDToPotentialActions[wRoll.ID]
		for _, action := range potentialActions {
			if state.ScoreThisRound == math.MaxUint8 && action.ContinueRolling {
				// Overflowed score this round. Our assumption is that this is unlikely.
				// Approximate the solution using the probability as if they stopped.
				action.ContinueRolling = false
			}

			newState := ApplyAction(state, action)
			if notYetOnBoard && !action.ContinueRolling && newState.PlayerScores[state.NumPlayers-1] < 500/incr {
				// Not a valid state: You must get at least 500 to get on the board.
				continue
			}

			if !recursiveEnumerateStates(newState, mask, depth+1, yield) {
				return false
			}
		}

		if len(potentialActions) == 0 {
			newState := ApplyAction(state, Action{})
			if !recursiveEnumerateStates(newState, mask, depth+1, yield) {
				return false
			}
		}
	}

	return yield(depth, state)
}

func init() {
	if scoreCache[0] != 0 {
		panic(fmt.Errorf("farkle should have zero score! got %d", scoreCache[0]))
	}
}