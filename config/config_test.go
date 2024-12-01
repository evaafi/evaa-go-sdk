package config

import (
	"testing"
)

func TestAssetConfig_GetJettonWalletAddress(t *testing.T) {
	t.Run("MainMainnetConfigAssets", func(t *testing.T) {
		cfg := GetMainMainnetConfig()

		if jettonAddress, err := cfg.Assets[USDT.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get USDT jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQD_kMQkK-A9-CQu3CdOnQUDZ2_8bY8Zrh1PvtE3hZpxvdRH" {
			t.Errorf("failed to get USDT jettonWalletAddress, got %s", jettonAddress.String())
		}

		if jettonAddress, err := cfg.Assets[STTON.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get STTON jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQBOn-1b-315ogeCl5lfPYW0ut6sjA2eq4LTdRv5vJJ1SsxX" {
			t.Errorf("failed to get STTON jettonWalletAddress, got %s", jettonAddress.String())
		}

		if jettonAddress, err := cfg.Assets[TSTON.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get TSTON jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQCRC0telhv1QESvTx24nNqUWB72zCysXQ0Bx97lVzucQ3Gr" {
			t.Errorf("failed to get TSTON jettonWalletAddress, got %s", jettonAddress.String())
		}

		if jettonAddress, err := cfg.Assets[JUSDT.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get JUSDT jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQBwbF0otxLyA3VrRjjC1q7i3G7LtoEpdyBjZEuNtrhC4drm" {
			t.Errorf("failed to get JUSDT jettonWalletAddress, got %s", jettonAddress.String())
		}

		if jettonAddress, err := cfg.Assets[JUSDC.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get JUSDC jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQCEHZi-CLX2ghdsNbY35umR1OFODG5ySHrtK6GItMdWv7dS" {
			t.Errorf("failed to get JUSDC jettonWalletAddress, got %s", jettonAddress.String())
		}
	})

	t.Run("LpMainnetConfigAssets", func(t *testing.T) {
		cfg := GetLpMainnetConfig()

		if jettonAddress, err := cfg.Assets[USDT.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get USDT jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQBvqcFlLgTUnjQdeuUt8DOqX6a1By3af5wD-0s6YuKLnKQs" {
			t.Errorf("failed to get USDT jettonWalletAddress, got %s", jettonAddress.String())
		}

		if jettonAddress, err := cfg.Assets[TON_STORM.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get TON_STORM jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQCmyubLpmAmclAA0j6qxkxbToA9z3vai97Cw5bA-tR0lIhb" {
			t.Errorf("failed to get TON_STORM jettonWalletAddress, got %s", jettonAddress.String())
		}

		if jettonAddress, err := cfg.Assets[USDT_STORM.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get USDT_STORM jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQC2ls4q8_NHtGHbAasDjb6ipJMcI0JvagynO4n0k4uUCUbi" {
			t.Errorf("failed to get USDT_STORM jettonWalletAddress, got %s", jettonAddress.String())
		}

		//if jettonAddress, err := cfg.Assets[TONUSDT_STONFI.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
		//	t.Errorf("failed to get TONUSDT_STONFI jettonWalletAddress, err: %s", err)
		//} else if jettonAddress.String() != "EQBwbF0otxLyA3VrRjjC1q7i3G7LtoEpdyBjZEuNtrhC4drm" {
		//	t.Errorf("failed to get TONUSDT_STONFI jettonWalletAddress, got %s", jettonAddress.String())
		//}

		if jettonAddress, err := cfg.Assets[TONUSDT_DEDUST.ID()].GetJettonWalletAddress(cfg.MasterAddress); err != nil {
			t.Errorf("failed to get TONUSDT_DEDUST jettonWalletAddress, err: %s", err)
		} else if jettonAddress.String() != "EQAXPGjhVFddTZeLDvE14TtYswhwaZJ64P-8EWmRWchcYpI6" {
			t.Errorf("failed to get TONUSDT_DEDUST jettonWalletAddress, got %s", jettonAddress.String())
		}
	})
}
