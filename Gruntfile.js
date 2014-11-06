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
		shell: {
			buildAndRun: {
				command: 'go build'
			}
		},
		watch: {
			stylesheets: {
				files: ['app/assets/stylesheets/**/*.less'],
				tasks: ['less'],
				options: { livereload: true }
			},
			scripts: {
				files: ['app/assets/javascripts/**/*.js'],
				tasks: ['concat'],
				options: { livereload: true }
			},
			app: {
				files: ['*.go'],
				tasks: ['shell:buildAndRun'],
				options: { livereload: true }
			}
		},
	})

	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-less')
	grunt.loadNpmTasks('grunt-contrib-watch')
	grunt.loadNpmTasks('grunt-shell')

	grunt.registerTask('default', ['less', 'concat', 'shell:buildAndRun'])
}