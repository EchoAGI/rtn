/*jshint -W030 */
({
    logLevel: 2,
	baseUrl: '../static/js',
	mainConfigFile: '../static/js/main.js',
	optimize: 'uglify2',
	uglify2: {
		output: {
			beautify: false
		},
		compress: {
			sequences: true,
			global_defs: {
				DEBUG: false
			}
		},
		warnings: false,
		mangle: true
	},
	wrap: false,
	useStrict: false,
	dir: './out',
	skipDirOptimize: true,
	removeCombined: true,
	modules: [
		{
			name: 'main',
			exclude: [
				'base'
			]
		},
		{
			name: 'base'
		},
		{
			name: 'app',
			exclude: [
				'main',
				'base'
			],
			inlineText: true,
		},
		{
			name: 'libs/pdf/pdf',
			dir: './out/libs/pdf',
			override: {
				skipModuleInsertion: true
			}
		},
		{
			name: 'libs/pdf/compatibility',
			dir: './out/libs/compatibility',
			override: {
				skipModuleInsertion: true
			}
		},
		{
			name: 'libs/pdf/pdf.worker',
			dir: './out/libs/pdf',
			override: {
				skipModuleInsertion: true
			}
		},
		{
			name: 'sandboxes/youtube',
			dir: './out/sandboxes',
			override: {
				skipModuleInsertion: true
			}
		},
		{
			name: 'sandboxes/pdf',
			dir: './out/sandboxes',
			override: {
				skipModuleInsertion: true
			}
		},
		{
			name: 'sandboxes/webodf',
			dir: './out/sandboxes',
			override: {
				skipModuleInsertion: true
			}
		}
	]
})
