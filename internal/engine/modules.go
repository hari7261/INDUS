package engine

func defaultModuleFactories() map[string]moduleFactory {
	return map[string]moduleFactory{
		"core":        func(e *Engine) Module { return &coreModule{engine: e} },
		"developer":   func(e *Engine) Module { return &developerModule{engine: e} },
		"environment": func(e *Engine) Module { return &environmentModule{engine: e} },
		"filesystem":  func(e *Engine) Module { return &filesystemModule{engine: e} },
		"network":     func(e *Engine) Module { return &networkModule{engine: e} },
		"package":     func(e *Engine) Module { return &packageModule{engine: e} },
		"project":     func(e *Engine) Module { return &projectModule{engine: e} },
		"system":      func(e *Engine) Module { return &systemModule{engine: e} },
		"task":        func(e *Engine) Module { return &taskModule{engine: e} },
		"terminal":    func(e *Engine) Module { return &terminalModule{engine: e} },
		"toolchain":   func(e *Engine) Module { return &toolchainModule{engine: e} },
		"update":      func(e *Engine) Module { return &updateModule{engine: e} },
		"workspace":   func(e *Engine) Module { return &workspaceModule{engine: e} },
	}
}
