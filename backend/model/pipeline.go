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

type PipelineSchedule struct {
	gorm.Model
	PipelineID      uint
	Pipeline        Pipeline
	UniqueOcurrence time.Time `json:"uniqueOccurrence"`
	CronExpression  string    `json:"cronExpression"`
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

type PipelineScheduleReq struct {
	ID              uint      `json:"id"`
	UniqueOcurrence time.Time `json:"uniqueOccurrence"`
	CronExpression  string    `json:"cronExpression"`
}

type StepDataNameAndType struct {
	Name        string `json:"name"`
	IsFirstStep bool   `json:"isFirstStep"`
	// Scripts
	ScriptType string `json:"scriptType"`
	// Scikit Datasets
	Dataset string `json:"dataset"`
}

type StepDataConfig struct {
	// CheckoutRepo
	RepoURL null.String `json:"repoURL"`
	// Scripts
	InlineScript null.String `json:"script"`
	Filename     null.String `json:"filename"`
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
	Max_iter            null.String `json:"max_iter"`
	Tol                 null.String `json:"tol"`
	Solver              null.String `json:"solver"`
	Random_state        null.String `json:"random_state"`
	Alphas              null.String `json:"alphas"`
	Class_weight        null.String `json:"class_weight"`
	Scoring             null.String `json:"scoring"`
	Cv                  null.String `json:"cv"`
	Store_cv_values     null.Bool   `json:"store_cv_values"`
	Precompute          null.Bool   `json:"precompute"`
	Warm_start          null.String `json:"warm_start"`
	Selection           null.String `json:"selection"`
	Eps                 null.String `json:"eps"`
	N_alphas            null.String `json:"n_alphas"`
	Verbose             null.Bool   `json:"verbose"`
	Fit_path            null.Bool   `json:"fit_path"`
	Jitter              null.String `json:"jitter"`
	Max_n_alphas        null.String `json:"max_n_alphas"`
	Criterion           null.String `json:"criterion"`
	Noise_variance      null.String `json:"noise_variance"`
	L1_ratio            null.String `json:"l1_ratio"`
	N_nonzero_coefs     null.String `json:"n_nonzero_coefs"`
	Copy                null.Bool   `json:"copy"`
	N_iter              null.String `json:"n_iter"`
	Alpha_1             null.String `json:"alpha_1"`
	Alpha_2             null.String `json:"alpha_2"`
	Lambda_1            null.String `json:"lambda_1"`
	Lambda_2            null.String `json:"lambda_2"`
	Alpha_init          null.String `json:"alpha_init"`
	Lambda_init         null.String `json:"lambda_init"`
	Compute_score       null.Bool   `json:"compute_score"`
	Threshold_lambda    null.String `json:"threshold_lambda"`
	Penalty             null.String `json:"penalty"`
	Dual                null.Bool   `json:"dual"`
	C                   null.String `json:"C"`
	Intercept_scaling   null.String `json:"intercept_scaling"`
	Multi_class         null.String `json:"multi_class"`
	Cs                  null.String `json:"Cs"`
	Refit               null.Bool   `json:"refit"`
	L1_ratios           null.String `json:"l1_ratios"`
	Power               null.String `json:"power"`
	Link                null.String `json:"link"`
	Loss                null.String `json:"loss"`
	Shuffle             null.Bool   `json:"shuffle"`
	Epsilon             null.String `json:"epsilon"`
	Learning_rate       null.String `json:"learning_rate"`
	Eta0                null.String `json:"eta0"`
	Power_t             null.String `json:"power_t"`
	Early_stopping      null.Bool   `json:"early_stopping"`
	Validation_fraction null.String `json:"validation_fraction"`
	N_iter_no_change    null.String `json:"n_iter_no_change"`
	Average             null.Bool   `json:"average"`
	Max_subpopulation   null.String `json:"max_subpopulation"`
	N_subsamples        null.String `json:"n_subsamples"`
	Quantile            null.String `json:"quantile"`
	Solver_options      null.String `json:"solver_options"`
	// HITL
	Data_dir         null.String `json:"data_dir"`
	Models_dir       null.String `json:"models_dir"`
	Epochs_dir       null.String `json:"epochs_dir"`
	Epochs           null.Int    `json:"epochs"`
	Tr_fraction      null.String `json:"tr_fraction"`
	Val_fraction     null.String `json:"val_fraction"`
	Train_desc       null.String `json:"train_desc"`
	Sampling         null.String `json:"sampling"`
	Entropy_thresh   null.String `json:"entropy_thresh"`
	Nr_queries       null.Int    `json:"nr_queries"`
	IsOversampled    null.Bool   `json:"isOversampled"`
	Start_epoch      null.Int    `json:"start_epoch"`
	Dataset          null.String `json:"dataset"`
	Pretrained_model null.String `json:"pretrained_model"`
	// Custom
	CustomArguments null.String `json:"customArguments"`
	// Dataset
	DatasetID   uint   `json:"datasetID"`
	DatasetName string `json:"datasetName"`
	DatasetPath string `json:"datasetPath"`
	// Trainer
	TrainerID   uint   `json:"trainerID"`
	TrainerName string `json:"trainerName"`
	TrainerPath string `json:"trainerPath"`
	IsStaggered bool   `json:"isStaggered"`
	// Tester
	TesterID   uint   `json:"testerID"`
	TesterName string `json:"testerName"`
	TesterPath string `json:"testerPath"`
	// Trained
	TrainedID   uint   `json:"trainedID"`
	TrainedName string `json:"trainedName"`
	TrainedPath string `json:"trainedPath"`
}

type StepData struct {
	ID          string              `json:"id"`
	NameAndType StepDataNameAndType `json:"nameAndType"`
	StepConfig  StepDataConfig      `json:"stepConfig"`
	Type        string              `json:"type"`
	// Scikit Unsupervised Models
	Model       null.String `json:"model"`
	IsFirstStep bool        `json:"isFirstStep"`
}

type NodeDescription struct {
	ID       string   `json:"id"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	SourceID string   `json:"source"`
	TargetID string   `json:"target"`
	Data     StepData `json:"data"`
}
