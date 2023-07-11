package steps

import (
	"di/model"
	"di/util"
	"di/util/errors"
	"os"
)

var ScikitUnsupervisedModelTypes = []string{
	"leastSquares",
	"ridgeRegression",
	"ridgeRegressionCV",
	"ridgeClassifier",
	"ridgeClassifierCV",
	"lasso",
	"lassoCV",
	"lassoLars",
	"lassoLarsCV",
	"lassoLarsIC",
	"multiTaskLasso",
	"multiTaskLassoCV",
	"elasticNet",
	"elasticNetCV",
	"multiTaskElasticNet",
	"multiTaskElasticNetCV",
	"lars",
	"larsCV",
	"omp",
	"ompCV",
	"bayesianRidge",
	"bayesianARD",
	"logisticRegression",
	"logisticRegressionCV",
	"tweedieRegressor",
	"poissonRegressor",
	"gammaRegressor",
	"sgdClassifier",
	"sgdRegressor",
	"perceptron",
	"passiveAgressiveClassifier",
	"passiveAgressiveRegressor",
	"huberRegression",
	"ransacRegression",
	"theilSenRegression",
	"quantileRegression",
}

type ScikitUnsupervisedModel struct {
	ID         uint
	PipelineID uint
	RunID      uint
	Model      string
	Name       string
	DataConfig model.StepDataConfig
}

func (step ScikitUnsupervisedModel) GetID() int {
	return int(step.ID)
}

func (step ScikitUnsupervisedModel) GetName() string {
	return step.Name
}

func (step *ScikitUnsupervisedModel) SetModel(model string) error {
	if util.StringArrayContains(ScikitUnsupervisedModelTypes, model) {
		step.Model = model
		return nil
	}

	return errors.NewNotFound("model", model)
}

func (step *ScikitUnsupervisedModel) SetConfig(stepConfig model.StepDataConfig) error {
	step.DataConfig = stepConfig
	return nil
}

func (step *ScikitUnsupervisedModel) SetPipelineID(pipelineID uint) error {
	step.PipelineID = pipelineID
	return nil
}

func (step *ScikitUnsupervisedModel) SetRunID(runID uint) error {
	step.RunID = runID

	return nil
}

func (step *ScikitUnsupervisedModel) GetPipelineID() uint {
	return step.PipelineID
}

func (step *ScikitUnsupervisedModel) GetRunID() uint {
	return step.RunID
}

func (step ScikitUnsupervisedModel) Execute(logFile *os.File) error {

	var args []string

	if step.DataConfig.Fit_intercept.Valid {
		args = append(args, "--fit_intercept")
	}

	if step.DataConfig.Copy_X.Valid && step.DataConfig.Copy_X.Bool {
		args = append(args, "--copy_X")
	}

	if step.DataConfig.N_jobs.Valid {
		args = append(args, "--n_jobs")
		args = append(args, string(step.DataConfig.N_jobs.Int64))
	}

	if step.DataConfig.Positive.Valid && step.DataConfig.Positive.Bool {
		args = append(args, "--positive")
	}

	if step.DataConfig.Alpha.Valid {
		args = append(args, "--alpha")
		args = append(args, step.DataConfig.Alpha.String)
	}

	if step.DataConfig.Max_iter.Valid {
		args = append(args, "--max_iter")
		args = append(args, string(step.DataConfig.Max_iter))
	}

	if step.DataConfig.Tol.Valid {
		args = append(args, "--tol")
		args = append(args, string(step.DataConfig.Tol))
	}

	if step.DataConfig.Solver.Valid {
		args = append(args, "--solver")
		args = append(args, string(step.DataConfig.Solver))
	}

	if step.DataConfig.Random_state.Valid {
		args = append(args, "--random_state")
		args = append(args, string(step.DataConfig.Random_state))
	}

	if step.DataConfig.Alphas.Valid {
		args = append(args, "--alphas")
		args = append(args, string(step.DataConfig.Alphas))
	}

	if step.DataConfig.Class_weight.Valid {
		args = append(args, "--class_weight")
		args = append(args, string(step.DataConfig.Class_weight))
	}

	if step.DataConfig.Scoring.Valid {
		args = append(args, "--scoring")
		args = append(args, string(step.DataConfig.Scoring))
	}

	if step.DataConfig.Cv.Valid {
		args = append(args, "--cv")
		args = append(args, string(step.DataConfig.Cv))
	}

	if step.DataConfig.Store_cv_values.Valid {
		args = append(args, "--store_cv_values")
		args = append(args, string(step.DataConfig.Store_cv_values))
	}

	if step.DataConfig.Precompute.Valid {
		args = append(args, "--precompute")
		args = append(args, string(step.DataConfig.Precompute))
	}

	if step.DataConfig.Warm_start.Valid {
		args = append(args, "--warm_start")
		args = append(args, string(step.DataConfig.Warm_start))
	}

	if step.DataConfig.Selection.Valid {
		args = append(args, "--selection")
		args = append(args, string(step.DataConfig.Selection))
	}

	if step.DataConfig.Eps.Valid {
		args = append(args, "--eps")
		args = append(args, string(step.DataConfig.Eps))
	}

	if step.DataConfig.N_alphas.Valid {
		args = append(args, "--n_alphas")
		args = append(args, string(step.DataConfig.N_alphas))
	}

	if step.DataConfig.Verbose.Valid {
		args = append(args, "--verbose")
		args = append(args, string(step.DataConfig.Verbose))
	}

	if step.DataConfig.Fit_path.Valid {
		args = append(args, "--fit_path")
		args = append(args, string(step.DataConfig.Fit_path))
	}

	if step.DataConfig.Jitter.Valid {
		args = append(args, "--jitter")
		args = append(args, string(step.DataConfig.Jitter))
	}

	if step.DataConfig.Max_n_alphas.Valid {
		args = append(args, "--max_n_alphas")
		args = append(args, string(step.DataConfig.Max_n_alphas))
	}

	if step.DataConfig.Criterion.Valid {
		args = append(args, "--criterion")
		args = append(args, string(step.DataConfig.Criterion))
	}

	if step.DataConfig.Noise_variance.Valid {
		args = append(args, "--noise_variance")
		args = append(args, string(step.DataConfig.Noise_variance))
	}

	if step.DataConfig.L1_ratio.Valid {
		args = append(args, "--l1_ratio")
		args = append(args, string(step.DataConfig.L1_ratio))
	}

	if step.DataConfig.N_nonzero_coefs.Valid {
		args = append(args, "--n_nonzero_coefs")
		args = append(args, string(step.DataConfig.N_nonzero_coefs))
	}

	if step.DataConfig.Copy.Valid {
		args = append(args, "--copy")
		args = append(args, string(step.DataConfig.Copy))
	}

	if step.DataConfig.N_iter.Valid {
		args = append(args, "--n_iter")
		args = append(args, string(step.DataConfig.N_iter))
	}

	if step.DataConfig.Alpha_1.Valid {
		args = append(args, "--alpha_1")
		args = append(args, string(step.DataConfig.Alpha_1))
	}

	if step.DataConfig.Alpha_2.Valid {
		args = append(args, "--alpha_2")
		args = append(args, string(step.DataConfig.Alpha_2))
	}

	if step.DataConfig.Lambda_1.Valid {
		args = append(args, "--lambda_1")
		args = append(args, string(step.DataConfig.Lambda_1))
	}

	if step.DataConfig.Lambda_2.Valid {
		args = append(args, "--lambda_2")
		args = append(args, string(step.DataConfig.Lambda_2))
	}

	if step.DataConfig.Alpha_init.Valid {
		args = append(args, "--alpha_init")
		args = append(args, string(step.DataConfig.Alpha_init))
	}

	if step.DataConfig.Lambda_init.Valid {
		args = append(args, "--lambda_init")
		args = append(args, string(step.DataConfig.Lambda_init))
	}

	if step.DataConfig.Compute_score.Valid {
		args = append(args, "--compute_score")
		args = append(args, string(step.DataConfig.Compute_score))
	}

	if step.DataConfig.Threshold_lambda.Valid {
		args = append(args, "--threshold_lambda")
		args = append(args, string(step.DataConfig.Threshold_lambda))
	}

	if step.DataConfig.Penalty.Valid {
		args = append(args, "--penalty")
		args = append(args, string(step.DataConfig.Penalty))
	}

	if step.DataConfig.Dual.Valid {
		args = append(args, "--dual")
		args = append(args, string(step.DataConfig.Dual))
	}

	if step.DataConfig.C.Valid {
		args = append(args, "--C")
		args = append(args, string(step.DataConfig.C))
	}

	if step.DataConfig.Intercept_scaling.Valid {
		args = append(args, "--intercept_scaling")
		args = append(args, string(step.DataConfig.Intercept_scaling))
	}

	if step.DataConfig.Multi_class.Valid {
		args = append(args, "--multi_class")
		args = append(args, string(step.DataConfig.Multi_class))
	}

	if step.DataConfig.Cs.Valid {
		args = append(args, "--Cs")
		args = append(args, string(step.DataConfig.Cs))
	}

	if step.DataConfig.Refit.Valid {
		args = append(args, "--refit")
		args = append(args, string(step.DataConfig.Refit))
	}

	if step.DataConfig.L1_ratios.Valid {
		args = append(args, "--l1_ratios")
		args = append(args, string(step.DataConfig.L1_ratios))
	}

	if step.DataConfig.Power.Valid {
		args = append(args, "--power")
		args = append(args, string(step.DataConfig.Power))
	}

	if step.DataConfig.Link.Valid {
		args = append(args, "--link")
		args = append(args, string(step.DataConfig.Link))
	}

	if step.DataConfig.Loss.Valid {
		args = append(args, "--loss")
		args = append(args, string(step.DataConfig.Loss))
	}

	if step.DataConfig.Shuffle.Valid {
		args = append(args, "--shuffle")
		args = append(args, string(step.DataConfig.Shuffle))
	}

	if step.DataConfig.Epsilon.Valid {
		args = append(args, "--epsilon")
		args = append(args, string(step.DataConfig.Epsilon))
	}

	if step.DataConfig.Learning_rate.Valid {
		args = append(args, "--learning_rate")
		args = append(args, string(step.DataConfig.Learning_rate))
	}

	if step.DataConfig.Eta0.Valid {
		args = append(args, "--eta0")
		args = append(args, string(step.DataConfig.Eta0))
	}

	if step.DataConfig.Power_t.Valid {
		args = append(args, "--power_t")
		args = append(args, string(step.DataConfig.Power_t))
	}

	if step.DataConfig.Early_stopping.Valid {
		args = append(args, "--early_stopping")
		args = append(args, string(step.DataConfig.Early_stopping))
	}

	if step.DataConfig.Validation_fraction.Valid {
		args = append(args, "--validation_fraction")
		args = append(args, string(step.DataConfig.Validation_fraction))
	}

	if step.DataConfig.N_iter_no_change.Valid {
		args = append(args, "--n_iter_no_change")
		args = append(args, string(step.DataConfig.N_iter_no_change))
	}

	if step.DataConfig.Average.Valid {
		args = append(args, "--average")
		args = append(args, string(step.DataConfig.Average))
	}

	if step.DataConfig.Max_subpopulation.Valid {
		args = append(args, "--max_subpopulation")
		args = append(args, string(step.DataConfig.Max_subpopulation))
	}

	if step.DataConfig.N_subsamples.Valid {
		args = append(args, "--n_subsamples")
		args = append(args, string(step.DataConfig.N_subsamples))
	}

	if step.DataConfig.Quantile.Valid {
		args = append(args, "--quantile")
		args = append(args, string(step.DataConfig.Quantile))
	}

	if step.DataConfig.Solver_options.Valid {
		args = append(args, "--solver_options")
		args = append(args, string(step.DataConfig.Solver_options))
	}

	return nil
}
