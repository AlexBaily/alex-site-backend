pipeline {
    agent any
    stages {
        stage('Test') {
            steps {
                script {

                    def siteImage = docker.build("alexsite-back:${env.BUILD_ID}")
                    siteImage.inside {
                        sh 'echo "Inside the container"'
                    }
                    siteImage.push('latest')
                } 
            }
        }
    }
}
