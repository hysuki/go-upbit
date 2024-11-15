package exchange

import (
	"encoding/json"
	"errors"
)

// MemberLevel은 사용자의 보안등급 정보를 나타내는 구조체입니다.
type MemberLevel struct {
	// SecurityLevel은 사용자의 보안등급입니다.
	SecurityLevel int `json:"security_level"`
	// FeeLevel은 사용자의 수수료등급입니다.
	FeeLevel int `json:"fee_level"`
	// EmailVerified는 사용자의 이메일 인증 여부입니다.
	EmailVerified bool `json:"email_verified"`
	// IdentityAuthVerified는 사용자의 실명 인증 여부입니다.
	IdentityAuthVerified bool `json:"identity_auth_verified"`
	// BankAccountVerified는 사용자의 계좌 인증 여부입니다.
	BankAccountVerified bool `json:"bank_account_verified"`
	// TwoFactorAuthVerified는 2FA 인증 수단의 활성화 여부입니다.
	TwoFactorAuthVerified bool `json:"two_factor_auth_verified"`
	// Locked는 사용자의 계정 보호 상태입니다.
	Locked bool `json:"locked"`
	// WalletLocked는 사용자의 출금 보호 상태입니다.
	WalletLocked bool `json:"wallet_locked"`
}

// WithdrawCurrency는 화폐 정보를 나타내는 구조체입니다.
type WithdrawCurrency struct {
	// Code는 화폐를 의미하는 영문 대문자 코드입니다.
	Code string `json:"code"`
	// WithdrawFee는 해당 화폐의 출금 수수료입니다.
	WithdrawFee string `json:"withdraw_fee"`
	// IsCoin은 화폐의 디지털 자산 여부입니다.
	IsCoin bool `json:"is_coin"`
	// WalletState는 해당 화폐의 지갑 상태입니다.
	WalletState string `json:"wallet_state"`
	// WalletSupport는 해당 화폐가 지원하는 입출금 정보입니다.
	WalletSupport []string `json:"wallet_support"`
}

// WithdrawLimit는 출금 제약 정보를 나타내는 구조체입니다.
type WithdrawLimit struct {
	// Currency는 화폐를 의미하는 영문 대문자 코드입니다.
	Currency string `json:"currency"`
	// Minimum은 출금 최소 금액/수량입니다.
	Minimum string `json:"minimum"`
	// RemainingDailyFiat은 통합 1일 잔여 출금 한도입니다.
	RemainingDailyFiat string `json:"remaining_daily_fiat"`
	// FiatCurrency는 해당 자산을 구매할 수 있는 법정 화폐입니다.
	FiatCurrency string `json:"fiat_currency"`
	// WithdrawDelayedFiat은 24시간 출금지연제에 의해 출금 제한되어 있는 금액입니다.
	WithdrawDelayedFiat string `json:"withdraw_delayed_fiat"`
	// Fixed는 출금 금액/수량 소수점 자리 수입니다.
	Fixed int `json:"fixed"`
	// CanWithdraw는 출금 지원 여부입니다.
	CanWithdraw bool `json:"can_withdraw"`
}

// WithdrawChance는 출금 가능 정보를 나타내는 구조체입니다.
type WithdrawChance struct {
	// MemberLevel은 사용자의 보안등급 정보입니다.
	MemberLevel MemberLevel `json:"member_level"`
	// Currency는 화폐 정보입니다.
	Currency WithdrawCurrency `json:"currency"`
	// Account는 사용자의 계좌 정보입니다.
	Account Accounts `json:"account"`
	// WithdrawLimit는 출금 제약 정보입니다.
	WithdrawLimit WithdrawLimit `json:"withdraw_limit"`
}

// GetWithdrawChance는 해당 통화의 가능한 출금 정보를 조회합니다.
// 엔드포인트: https://api.upbit.com/v1/withdraws/chance
func (e *Exchange) GetWithdrawChance(currency string) (*WithdrawChance, error) {
	if currency == "" {
		return nil, errors.New("currency is required")
	}

	params := map[string]string{
		"currency": currency,
	}

	resp, err := e.Client.Get("/withdraws/chance", params)
	if err != nil {
		return nil, err
	}

	var chance WithdrawChance
	if err := json.Unmarshal(resp, &chance); err != nil {
		return nil, err
	}

	return &chance, nil
}
