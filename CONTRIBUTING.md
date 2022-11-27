# Contributing
* Squash your commits.
* Please open an issue first.
* _"If you like it, then you shoulda put a TEST on it."_ - Beyoncé

# Adding a new secrets management system
Adding a support for another system is as easy as filling the `SecretValues` struct and giving at as a value to `kubesealconvert.BuildSecretFile(secretData)`.  
`BuildSecretFile` will transform it into a normal k8s secret, and will run the `kubeseal` command.