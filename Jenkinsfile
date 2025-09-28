pipeline {
    agent any

    environment {
        CGO_ENABLED = '1'
        ORAPIKEY = credentials('openrouter-api-key')
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }

        stage('Setup Go') {
            steps {
                script {
                    def goVersion = '1.23'
                    sh """
                        if ! command -v go &> /dev/null; then
                            echo "Installing Go ${goVersion}..."
                            wget -q https://go.dev/dl/go${goVersion}.linux-amd64.tar.gz
                            sudo tar -C /usr/local -xzf go${goVersion}.linux-amd64.tar.gz
                            export PATH=\$PATH:/usr/local/go/bin
                        fi
                        go version
                    """
                }
            }
        }

        stage('Dependencies') {
            steps {
                sh 'go mod download'
                sh 'go mod tidy'
            }
        }

        stage('Unit Tests') {
            steps {
                sh 'go test ./...'
            }
        }

        stage('Coverage Tests') {
            steps {
                sh 'go test -cover ./...'
            }
        }

        stage('Race Detection') {
            steps {
                sh 'go test -race ./...'
            }
        }

        stage('Specific Tests') {
            steps {
                sh 'go test -run TestChatComplete'
            }
        }

        stage('Integration Test') {
            steps {
                sh 'go run cmd/openrouter-test/main.go -key $ORAPIKEY -test all -model "google/gemini-2.5-flash-lite"'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Test suite completed successfully!'
        }
        failure {
            echo 'Test suite failed. Please check the logs.'
        }
    }
}