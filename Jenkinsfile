import org.jenkinsci.plugins.pipeline.github.trigger.IssueCommentCause

@Library('pipeline-lib') _
@Library('cve-monitor') __

def MAIN_BRANCH                    = 'master'
def DOCKER_REPOSITORY_NAME         = 'cloudfront-signer'
def DOCKER_REGISTRY_URL            = 'https://662491802882.dkr.ecr.us-east-1.amazonaws.com'
def DOCKER_REGISTRY_CREDENTIALS_ID = 'ecr:us-east-1:ecr-docker-push'

properties([
    pipelineTriggers([issueCommentTrigger('!build')])
])
def isForcePublish = !!currentBuild.rawBuild.getCause(IssueCommentCause)

withResultReporting(slackChannel: '#tm-inf', mainBranch: MAIN_BRANCH) {
  inDockerAgent(deployer.wrapPodTemplate(containers: [imageScanner.container()])) {
    def version
    def dockerImage

    stage('Build docker image') {
      checkout([
        $class: 'GitSCM',
        branches: scm.branches,
        doGenerateSubmoduleConfigurations: scm.doGenerateSubmoduleConfigurations,
        extensions: scm.extensions + [[$class: 'CloneOption', noTags: false, shallow: false, depth: 0, reference: '']],
        userRemoteConfigs: scm.userRemoteConfigs,
      ])

      version = sh(returnStdout: true, script: 'git describe --tags --always --dirty').trim()
      dockerImage = docker.build(DOCKER_REPOSITORY_NAME)
    }
    stage('Scan image') {
      imageScanner.scan(dockerImage)
    }
    if (BRANCH_NAME == MAIN_BRANCH || isForcePublish) {
      stage('Publish docker image') {
        docker.withRegistry(DOCKER_REGISTRY_URL, DOCKER_REGISTRY_CREDENTIALS_ID) {
          echo("Publishing docker image ${dockerImage.imageName()} with tag ${version}")
          dockerImage.push(version)
          if (BRANCH_NAME == MAIN_BRANCH) {
            dockerImage.push("latest")
          }
        }
        if (isForcePublish) {
          pullRequest.comment("Built and published ${version}")
        }
      }
    }
  }
}
