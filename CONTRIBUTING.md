# Contributing
* Squash your commits.
* Please open an issue first.
* _"If you like it, then you shoulda put a TEST on it."_ - Beyoncé

# Adding a new secrets management system
Adding a support for another system is as easy as filling the `SecretValues` struct and giving at as a value to `KubeSeal.BuildSecretFile(secretData)`.  
`BuildSecretFile` will transform it into a normal k8s secret, and will run the `kubeseal` command.  
Make sure to create the proper interface for your system (might be extended later).  
Check out [Vault implementation](https://github.com/EladLeev/kubeseal-convert/pull/10/files) as an example!


## Contributors List
👏 @Pavel-Durov
