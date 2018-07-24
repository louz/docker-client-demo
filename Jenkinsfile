pipeline {
   environment {
     REGISTRY_ENDPOINT = "https://docker-client-demo/v2/"
     IMAGE_WITH_TAG = "docker-client-demo"
     REGISTRY_CERTS = "registry"
   }
  agent {
    node {
      label 'golang'
    }

  }
  stages {
    stage('Build') {
      steps {
        sh 'go build -o app'
      }
    }
    stage('Test') {
      steps {
        sh 'go test'
      }
    }
    stage('Code Quality') {
      steps {
        script {
          try{
            checkstyle canComputeNew: false, defaultEncoding: '', healthy: '', pattern: '', unHealthy: ''
          }catch(e){
            echo e
          }
        }

      }
    }
    stage('Image Build&Publish') {
      steps {
        echo 'Build Images'
        script {
          docker.withRegistry("${REGISTRY_ENDPOINT}", "${REGISTRY_CERTS}") {
            sh 'docker build -t ${IMAGE_WITH_TAG} .'
            sh 'docker push ${IMAGE_WITH_TAG}'
          }
        }

      }
    }
  }
  triggers {
    pollSCM('* * * * *')
  }
}