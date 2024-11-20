package provider

import (
	"strings"

	"go.uber.org/fx/fxevent"
)

// LogEvent logs the given event to the provided Zerolog.
func (z *ZeroLogger) FxLogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		z.Info(
			"OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)

	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			z.Error(
				"OnStart hook failed",
				"error", e.Err,
				"callee", e.FunctionName,
				"caller", e.CallerName,
			)
		} else {
			z.Info(
				"OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		z.Info(
			"OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			z.Error(
				"OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"error", e.Err,
			)
		} else {
			z.Info(
				"OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		z.Error(
			"supplied",
			"error", e.Err,
			"type", e.TypeName,
			"module", e.ModuleName,
		)
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			z.Info(
				"provided",
				"constructor", e.ConstructorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			z.Error(
				"error encountered while applying options",
				"error", e.Err,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			z.Info(
				"decorated",
				"decorator", e.DecoratorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			z.Error(
				"error encountered while applying options",
				"error", e.Err,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		z.Info(
			"invoking",
			"function", e.FunctionName,
			"module", e.ModuleName,
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			z.Error(
				"invoke failed",
				"error", e.Err,
				"stack", e.Trace,
				"function", e.FunctionName,
			)
		}
	case *fxevent.Stopping:
		z.Info(
			"received signal",
			"signal", strings.ToUpper(e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			z.Error(
				"stop failed",
				"error", e.Err,
			)
		}
	case *fxevent.RollingBack:
		z.Error(
			"start failed, rolling back",
			"error", e.StartErr,
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			z.Error(
				"rollback failed",
				"error", e.Err,
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			z.Error(
				"start failed",
				"error", e.Err,
			)
		} else {
			z.Info("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			z.Error(
				"custom logger initialization failed",
				"error", e.Err,
			)
		} else {
			z.Info(
				"initialized custom fxevent.Logger",
				"function", e.ConstructorName,
			)
		}
	}
}
