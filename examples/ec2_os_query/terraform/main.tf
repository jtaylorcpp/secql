provider "aws" {
    region = "us-west-1"
}

data "aws_security_group" "default" {
  name   = "default"
  vpc_id = module.vpc.vpc_id
}

module "vpc" {
  source = "github.com/terraform-aws-modules/terraform-aws-vpc"

  name = "simple-example"

  cidr = "10.0.0.0/16"

  azs             = ["us-west-1a", "us-west-1c"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_ipv6 = false

  enable_nat_gateway = true
  single_nat_gateway = true

  public_subnet_tags = {
    Name = "overridden-name-public"
  }

  tags = {
    Owner       = "secql"
    Environment = "example-ec2"
  }

  vpc_tags = {
    Name = "secql-example-ec2"
  }
}

locals {
  osquery_install = <<OSQ
#!/bin/bash
apt-get update -y && apt-get upgrade -y

export OSQUERY_KEY=1484120AC4E9F8A1A577AEEE97A80C63C9D8B80B
apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys $OSQUERY_KEY
add-apt-repository 'deb [arch=amd64] https://pkg.osquery.io/deb deb main'
apt-get update -y
apt-get install osquery -y

apt-get install wget -y
wget https://github.com/jtaylorcpp/secql/releases/download/v0.0.0/secqld_linux
chmod +x secqld_linux
./secqld_linux systemd install-secqld
systemctl enable secqld
systemctl restart secqld
  OSQ
}
data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "osquery" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"

  user_data = local.osquery_install

  subnet_id = module.vpc.public_subnets[0]
  associate_public_ip_address = true
  vpc_security_group_ids = [
    aws_security_group.basics.id
  ]

  tags = {
    Name = "secql-osquery"
  }
}

data "http" "myip" {
  url = "http://ipv4.icanhazip.com"
}

resource "aws_security_group" "basics" {
  name        = "basics"
  description = "Allow basic inbound traffic"
  vpc_id      = module.vpc.vpc_id

  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = [
      "${chomp(data.http.myip.body)}/32"
    ]
  }

  ingress {
    description = "OSQuery Agent"
    from_port   = 8000
    to_port     = 8000
    protocol    = "tcp"
    cidr_blocks = [
      "${chomp(data.http.myip.body)}/32"
    ]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "basics"
  }

  depends_on = [
    module.vpc
  ]
}

output "module" {
  value = module.vpc
}

output "ec2" {
  value = {
    id = aws_instance.osquery.id
  }
}