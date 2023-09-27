package entity

type TxResult struct {
	Hash    string `json:"hash"`
	Lt      int64  `json:"lt"`
	Account struct {
		Address string `json:"address"`
		IsScam  bool   `json:"is_scam"`
	} `json:"account"`
	Success         bool   `json:"success"`
	Utime           int64  `json:"utime"`
	OrigStatus      string `json:"orig_status"`
	EndStatus       string `json:"end_status"`
	TotalFees       int    `json:"total_fees"`
	TransactionType string `json:"transaction_type"`
	StateUpdateOld  string `json:"state_update_old"`
	StateUpdateNew  string `json:"state_update_new"`
	InMsg           struct {
		CreatedLt   int64 `json:"created_lt"`
		IhrDisabled bool  `json:"ihr_disabled"`
		Bounce      bool  `json:"bounce"`
		Bounced     bool  `json:"bounced"`
		Value       int   `json:"value"`
		FwdFee      int   `json:"fwd_fee"`
		IhrFee      int   `json:"ihr_fee"`
		Destination struct {
			Address string `json:"address"`
			IsScam  bool   `json:"is_scam"`
		} `json:"destination"`
		ImportFee int    `json:"import_fee"`
		CreatedAt int    `json:"created_at"`
		RawBody   string `json:"raw_body"`
	} `json:"in_msg"`
	OutMsgs []struct {
		CreatedLt   int64 `json:"created_lt"`
		IhrDisabled bool  `json:"ihr_disabled"`
		Bounce      bool  `json:"bounce"`
		Bounced     bool  `json:"bounced"`
		Value       int   `json:"value"`
		FwdFee      int   `json:"fwd_fee"`
		IhrFee      int   `json:"ihr_fee"`
		Destination struct {
			Address string `json:"address"`
			IsScam  bool   `json:"is_scam"`
		} `json:"destination"`
		Source struct {
			Address string `json:"address"`
			IsScam  bool   `json:"is_scam"`
		} `json:"source"`
		ImportFee     int    `json:"import_fee"`
		CreatedAt     int    `json:"created_at"`
		OpCode        string `json:"op_code"`
		RawBody       string `json:"raw_body"`
		DecodedOpName string `json:"decoded_op_name"`
		DecodedBody   struct {
			Text string `json:"text"`
		} `json:"decoded_body"`
	} `json:"out_msgs"`
	Block         string `json:"block"`
	PrevTransHash string `json:"prev_trans_hash"`
	PrevTransLt   int64  `json:"prev_trans_lt"`
	ComputePhase  struct {
		Skipped  bool `json:"skipped"`
		Success  bool `json:"success"`
		GasFees  int  `json:"gas_fees"`
		GasUsed  int  `json:"gas_used"`
		VmSteps  int  `json:"vm_steps"`
		ExitCode int  `json:"exit_code"`
	} `json:"compute_phase"`
	StoragePhase struct {
		FeesCollected int    `json:"fees_collected"`
		StatusChange  string `json:"status_change"`
	} `json:"storage_phase"`
	ActionPhase struct {
		Success        bool `json:"success"`
		TotalActions   int  `json:"total_actions"`
		SkippedActions int  `json:"skipped_actions"`
		FwdFees        int  `json:"fwd_fees"`
		TotalFees      int  `json:"total_fees"`
	} `json:"action_phase"`
	Aborted   bool `json:"aborted"`
	Destroyed bool `json:"destroyed"`
}
