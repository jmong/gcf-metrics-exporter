package metricsexporter
/**
 * Measures elapsed time.
 * 
 * @usage
 * t := NewTimer()
 * t.Start()
 * // Do some stuff
 * t.End()
 * fmt.Printf("%d elapsed", t.GetElapsed())
 **/

import (
)

/* @TODO
 */
type Timer struct {
    //start 
    //end
    tchan  chan struct{}
}

/* @TODO
 */
func NewTimer() *Timer {
    return &Timer{
        //tchan: make(chan struct{})
    }
}

/* @TODO
 */
func Start() {
}

/* @TODO
 */
func startTimer(tchan chan struct{}) {
}

/* @TODO
 */
func End() {
}

/* @TODO
 */
//func GetElapsed() int64 {
//}
