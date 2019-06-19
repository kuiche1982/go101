output "OutServiceName" {
  description = "The Service Name"
  value       = "${var.serviceName}"
}

output "IP" {
  description = "the eip ip"
  value       = "${aws_elb.example_elb.dns_name}"
}
