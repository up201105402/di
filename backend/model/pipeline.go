package model

import (
	"time"

	"gopkg.in/guregu/null.v4"
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
	RepoURL null.String `json:"repoURL"`
	// Scikit Datasets

	FilePath         null.String `json:"filePath"`
	DataFilePath     null.String `json:"dataFilePath"`
	TargetFilePath   null.String `json:"targetFilePath"`
	LowerXRangeIndex null.Int    `json:"lowerXRangeIndex"`
	UpperXRangeIndex null.Int    `json:"upperXRangeIndex"`
	LowerYRangeIndex null.Int    `json:"lowerYRangeIndex"`
	UpperYRangeIndex null.Int    `json:"upperYRangeIndex"`
	// Scikit Unsupervised Models
	Fit_intercept       null.Bool   `json:"fit_intercept"`
	Copy_X              null.Bool   `json:"copy_X"`
	N_jobs              null.Int    `json:"n_jobs"`
	Positive            null.Bool   `json:"positive"`
	Alpha               null.String `json:"alpha"`
	Max_iter            null.Float  `json:"max_iter"`
	Tol                 null.Float  `json:"tol"`
	Solver              null.String `json:"solver"`
	Random_state        null.Float  `json:"random_state"`
	Alphas              null.String `json:"alphas"`
	Class_weight        null.String `json:"class_weight"`
	Scoring             null.String `json:"scoring"`
	Cv                  null.Float  `json:"cv"`
	Store_cv_values     null.Bool   `json:"store_cv_values"`
	Precompute          null.Bool   `json:"precompute"`
	Warm_start          null.Float  `json:"warm_start"`
	Selection           null.String `json:"selection"`
	Eps                 null.Float  `json:"eps"`
	N_alphas            null.Float  `json:"n_alphas"`
	Verbose             null.Bool   `json:"verbose"`
	Fit_path            null.Bool   `json:"fit_path"`
	Jitter              null.Float  `json:"jitter"`
	Max_n_alphas        null.Float  `json:"max_n_alphas"`
	Criterion           null.String `json:"criterion"`
	Noise_variance      null.Float  `json:"noise_variance"`
	L1_ratio            null.Float  `json:"l1_ratio"`
	N_nonzero_coefs     null.Float  `json:"n_nonzero_coefs"`
	Copy                null.Bool   `json:"copy"`
	N_iter              null.Float  `json:"n_iter"`
	Alpha_1             null.Float  `json:"alpha_1"`
	Alpha_2             null.Float  `json:"alpha_2"`
	Lambda_1            null.Float  `json:"lambda_1"`
	Lambda_2            null.Float  `json:"lambda_2"`
	Alpha_init          null.Float  `json:"alpha_init"`
	Lambda_init         null.Float  `json:"lambda_init"`
	Compute_score       null.Bool   `json:"compute_score"`
	Threshold_lambda    null.Float  `json:"threshold_lambda"`
	Penalty             null.String `json:"penalty"`
	Dual                null.Bool   `json:"dual"`
	C                   null.Float  `json:"C"`
	Intercept_scaling   null.Float  `json:"intercept_scaling"`
	Multi_class         null.String `json:"multi_class"`
	Cs                  null.Float  `json:"Cs"`
	Refit               null.Bool   `json:"refit"`
	L1_ratios           null.Float  `json:"l1_ratios"`
	Power               null.Float  `json:"power"`
	Link                null.String `json:"link"`
	Loss                null.String `json:"loss"`
	Shuffle             null.Bool   `json:"shuffle"`
	Epsilon             null.Float  `json:"epsilon"`
	Learning_rate       null.String `json:"learning_rate"`
	Eta0                null.Float  `json:"eta0"`
	Power_t             null.Float  `json:"power_t"`
	Early_stopping      null.Bool   `json:"early_stopping"`
	Validation_fraction null.Float  `json:"validation_fraction"`
	N_iter_no_change    null.Float  `json:"n_iter_no_change"`
	Average             null.Bool   `json:"average"`
	Max_subpopulation   null.Float  `json:"max_subpopulation"`
	N_subsamples        null.Float  `json:"n_subsamples"`
	Quantile            null.Float  `json:"quantile"`
	Solver_options      null.String `json:"solver_options"`
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
