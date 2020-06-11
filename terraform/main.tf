terraform {
  backend "s3" {
    bucket     = "terraform-state-storage-586877430255"
    lock_table = "terraform-state-lock-586877430255"
    region     = "us-west-2"

    // THIS MUST BE UNIQUE
    key = "visca-service.tfstate"
  }
}

provider "aws" {
  region = "us-west-2"
}

data "aws_ssm_parameter" "eks_cluster_endpoint" {
  name = "/eks/av-cluster-endpoint"
}

provider "kubernetes" {
  host = data.aws_ssm_parameter.eks_cluster_endpoint.value
}

data "aws_ssm_parameter" "event_url" {
  name = "/env/visca-service/event-url"
}

data "aws_ssm_parameter" "dns_addr" {
  name = "/env/visca-service/dns-addr"
}

module "dev" {
  source = "github.com/byuoitav/terraform//modules/kubernetes-deployment"

  // required
  name           = "visca-service-dev"
  image          = "docker.pkg.github.com/byuoitav/visca-service/visca-service-dev"
  image_version  = "e10df7c"
  container_port = 8080
  repo_url       = "https://github.com/byuoitav/visca-service"

  // optional
  image_pull_secret = "github-docker-registry"
  public_urls       = ["visca-dev.av.byu.edu"]
  container_env     = {}
  container_args = [
    "--port", "8080",
    "--log-level", "0", // set log level to info
    "--name", "visca-service-dev",
    "--event-url", data.aws_ssm_parameter.event_url.value,
    "--dns-addr", data.aws_ssm_parameter.dns_addr.value,
  ]
  ingress_annotations = {
    // "nginx.ingress.kubernetes.io/whitelist-source-range" = "128.187.0.0/16"
  }
}
