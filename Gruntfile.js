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
		watch: {
			all: {
				files: ['app/assets/stylesheets/**/*.less', 'app/assets/javascripts/**/*.js'],
				tasks: ['default'],
				options: {
					livereload: true
				}
			}
		},
	})

	grunt.loadNpmTasks('grunt-contrib-concat');
	grunt.loadNpmTasks('grunt-contrib-less')
	grunt.loadNpmTasks('grunt-contrib-watch')

	grunt.registerTask('default', ['less', 'concat'])
}