package model

import (
	"time"

	"gorm.io/gorm"
)

type Pipeline struct {
	gorm.Model
	UserID     uint `json:"userId"`
	User       User
	Name       string    `json:"name"`
	Definition string    `json:"definition"`
	LastRun    time.Time `gorm:"-:all"`
}

type CreatePipelineReq struct {
	User       string `json:"user"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type PipelineReq struct {
	ID         uint   `json:"id"`
	User       string `json:"user"`
	Name       string `json:"name"`
	Definition string `json:"definition"`
}

type StepDataNameAndType struct {
	Name string `json:"name"`
	// Scikit Datasets
	Dataset     string `json:"dataset"`
	IsFirstStep bool   `json:"isFirstStep"`
}

type StepDataConfig struct {
	// CheckoutRepo
	RepoURL string `json:"repoURL"`
	// Scikit Datasets

	FilePath         string `json:"filePath"`
	DataFilePath     string `json:"dataFilePath"`
	TargetFilePath   string `json:"targetFilePath"`
	LowerXRangeIndex int    `json:"lowerXRangeIndex"`
	UpperXRangeIndex int    `json:"upperXRangeIndex"`
	LowerYRangeIndex int    `json:"lowerYRangeIndex"`
	UpperYRangeIndex int    `json:"upperYRangeIndex"`
	// Scikit Unsupervised Models
	Fit_intercept       bool    `json:"fit_intercept"`
	Copy_X              bool    `json:"copy_X"`
	N_jobs              float64 `json:"n_jobs"`
	Positive            bool    `json:"positive"`
	Alpha               float64 `json:"alpha"`
	Max_iter            float64 `json:"max_iter"`
	Tol                 float64 `json:"tol"`
	Solver              string  `json:"solver"`
	Random_state        float64 `json:"random_state"`
	Alphas              string  `json:"alphas"`
	Class_weight        string  `json:"class_weight"`
	Scoring             string  `json:"scoring"`
	Cv                  float64 `json:"cv"`
	Store_cv_values     bool    `json:"store_cv_values"`
	Precompute          bool    `json:"precompute"`
	Warm_start          float64 `json:"warm_start"`
	Selection           string  `json:"selection"`
	Eps                 float64 `json:"eps"`
	N_alphas            float64 `json:"n_alphas"`
	Verbose             bool    `json:"verbose"`
	Fit_path            bool    `json:"fit_path"`
	Jitter              float64 `json:"jitter"`
	Max_n_alphas        float64 `json:"max_n_alphas"`
	Criterion           string  `json:"criterion"`
	Noise_variance      float64 `json:"noise_variance"`
	L1_ratio            float64 `json:"l1_ratio"`
	N_nonzero_coefs     float64 `json:"n_nonzero_coefs"`
	Copy                bool    `json:"copy"`
	N_iter              float64 `json:"n_iter"`
	Alpha_1             float64 `json:"alpha_1"`
	Alpha_2             float64 `json:"alpha_2"`
	Lambda_1            float64 `json:"lambda_1"`
	Lambda_2            float64 `json:"lambda_2"`
	Alpha_init          float64 `json:"alpha_init"`
	Lambda_init         float64 `json:"lambda_init"`
	Compute_score       bool    `json:"compute_score"`
	Threshold_lambda    float64 `json:"threshold_lambda"`
	Penalty             string  `json:"penalty"`
	Dual                bool    `json:"dual"`
	C                   float64 `json:"C"`
	Intercept_scaling   float64 `json:"intercept_scaling"`
	Multi_class         string  `json:"multi_class"`
	Cs                  float64 `json:"Cs"`
	Refit               bool    `json:"refit"`
	L1_ratios           float64 `json:"l1_ratios"`
	Power               float64 `json:"power"`
	Link                string  `json:"link"`
	Loss                string  `json:"loss"`
	Shuffle             bool    `json:"shuffle"`
	Epsilon             float64 `json:"epsilon"`
	Learning_rate       string  `json:"learning_rate"`
	Eta0                float64 `json:"eta0"`
	Power_t             float64 `json:"power_t"`
	Early_stopping      bool    `json:"early_stopping"`
	Validation_fraction float64 `json:"validation_fraction"`
	N_iter_no_change    float64 `json:"n_iter_no_change"`
	Average             bool    `json:"average"`
	Max_subpopulation   float64 `json:"max_subpopulation"`
	N_subsamples        float64 `json:"n_subsamples"`
	Quantile            float64 `json:"quantile"`
	Solver_options      string  `json:"solver_options"`
}

type StepData struct {
	NameAndType StepDataNameAndType `json:"nameAndType"`
	StepConfig  StepDataConfig      `json:"stepConfig"`
	Type        string              `json:"type"`
	IsFirstStep bool                `json:"isFirstStep"`
}

type NodeDescription struct {
	ID       string   `json:"id"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	SourceID string   `json:"source"`
	TargetID string   `json:"target"`
	Data     StepData `json:"data"`
}
