provider "aws" {
  access_key = ""
  secret_key = ""
  region     = "ap-southeast-1"

  // access_key = "${var.access_key}"

  // secret_key = "${var.secret_key}"
  // region     = "${var.region}"
}

resource "aws_elb" "example_elb" {
  name = "${var.serviceName}-elb"

  # The same availability zone as our instances
  availability_zones = ["${aws_instance.example.*.availability_zone}"]

  listener {
    instance_port     = 80
    instance_protocol = "http"
    lb_port           = 80
    lb_protocol       = "http"
  }

  instances = ["${aws_instance.example.*.id}"]
}

resource "aws_instance" "example" {
  ami           = "ami-2757f631"
  instance_type = "t2.micro"

  count = 2

  tags = {
    T1 = "${var.serviceName}-web"
  }
}
