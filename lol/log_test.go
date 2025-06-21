package lol

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	log, chk, errorf := Main.Log, Main.Check, Main.Errorf
	for i := Off; i <= Trace; i++ {
		SetLoggers(i)
		fmt.Println(LevelNames[i])
		log.F.Ln("ln")
		log.E.Ln("ln")
		log.W.Ln("ln")
		log.I.Ln("ln")
		log.D.Ln("ln")
		log.T.Ln("ln")
		log.F.F("F: %s", LevelNames[i])
		log.E.F("F: %s", LevelNames[i])
		log.W.F("F: %s", LevelNames[i])
		log.I.F("F: %s", LevelNames[i])
		log.D.F("F: %s", LevelNames[i])
		log.T.F("F: %s", LevelNames[i])
		_ = errorf.F(LevelNames[i])
		_ = errorf.E(LevelNames[i])
		_ = errorf.W(LevelNames[i])
		_ = errorf.I(LevelNames[i])
		_ = errorf.D(LevelNames[i])
		_ = errorf.T(LevelNames[i])
		_ = chk.F(errorf.F(LevelNames[i]))
		_ = chk.E(errorf.E(LevelNames[i]))
		_ = chk.W(errorf.W(LevelNames[i]))
		_ = chk.I(errorf.I(LevelNames[i]))
		_ = chk.D(errorf.D(LevelNames[i]))
		_ = chk.T(errorf.T(LevelNames[i]))
		fmt.Println()
	}
	_, _ = chk, errorf
}
