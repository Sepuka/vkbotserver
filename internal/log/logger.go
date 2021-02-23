package log

import (
	"errors"
	errPkg "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(production bool) (*zap.SugaredLogger, error) {
	var (
		err               error
		logger            *zap.Logger
		sugar             *zap.SugaredLogger
		zapCfg            zap.Config
		core              zapcore.Core
		fileEncoder       zapcore.Encoder
		fileEncoderConfig zapcore.EncoderConfig
	)

	fileSynchronizer, closeOut, err := zap.Open(`stdout`)
	if err != nil {
		return nil, errPkg.Wrap(err, `unable to open output files`)
	}

	writeSyncer := zapcore.AddSync(fileSynchronizer)

	consoleMsgLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if production {
			return lvl >= zapcore.InfoLevel
		}

		return true
	})

	if production {
		zapCfg = zap.NewProductionConfig()
		fileEncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
		fileEncoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	zapCfg.OutputPaths = []string{`stdout`}

	fileEncoder = zapcore.NewJSONEncoder(fileEncoderConfig)
	core = zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writeSyncer, consoleMsgLevel),
	)

	logger = zap.New(core)
	sugar = logger.Sugar()
	if sugar == nil {
		closeOut()
		return nil, errors.New(`unable build sugar logger`)
	}

	return sugar, err
}
