module terraform-provider-secberus

go 1.13

replace github.com/RexBelli/go-secberus => /home/rex/go/src/github.com/RexBelli/go-secberus

require (
	github.com/RexBelli/go-secberus v0.0.0-20201119024355-eac22371d346
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.2.0
)