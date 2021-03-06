module.exports = function(grunt) {
	grunt.initConfig({
		pkg: grunt.file.readJSON('package.json'),
		less: {
			dist: {
				src: [
					'app/assets/stylesheets/app.less'
				],
				dest: 'public/app.css'
			}
		},
		cssmin: {
			minify: {
				src: ['public/app.css'],
				dest: 'public/app.css'
			}
		},
		concat: {
			dist: {
				src: [
					'app/assets/javascripts/jquery-2.1.1.min.js',
					'app/assets/javascripts/chat.js',
					'app/assets/javascripts/app.js'
				],
				dest: 'public/app.js'
			}
		},
		uglify: {
			my_target: {
				files: {
					'public/app.js': ['public/app.js']
				}
			}
		},
		shell: {
			buildAndRun: {
				command: 'go build'
			}
		},
		watch: {
			stylesheets: {
				files: ['app/assets/stylesheets/**/*.less'],
				tasks: ['less', 'cssmin'],
				options: { livereload: true }
			},
			scripts: {
				files: ['app/assets/javascripts/**/*.js'],
				tasks: ['concat', 'uglify'],
				options: { livereload: true }
			},
			app: {
				files: ['*.go'],
				tasks: ['shell:buildAndRun'],
				options: { livereload: true }
			}
		},
	})

	grunt.loadNpmTasks('grunt-contrib-concat')
	grunt.loadNpmTasks('grunt-contrib-cssmin')
	grunt.loadNpmTasks('grunt-contrib-less')
	grunt.loadNpmTasks('grunt-contrib-uglify')
	grunt.loadNpmTasks('grunt-contrib-watch')
	grunt.loadNpmTasks('grunt-shell')

	grunt.registerTask('default', ['shell:buildAndRun', 'less', 'cssmin', 'concat', 'uglify'])
}