# OCT: Open Container Testing

The oct(open container testing) project aims to promote the [Open Container Initiative](http://www.opencontainers.org/) by providing a universal testing `lib/tools` for all the container projects compliance with [oci specs](https://github.com/opencontainers/specs).

The testing tools are all maintained in the `tools` directory.

```
tools
  |______ octest
  |______ specsValidator

```

Now we have two tools `octest` and `specsValidator`.
You can join us to develop these tools or add new ones.

### [octest](tools/octest/README.md) - spec test (config and bundle)
- Verify whether a config file (config.json and runtime.json) is compliance with the newest OCI specs.
- Verify whether a bundle is valid according to the newest OCI specs.

### [specsValidator](tools/specsValidator/README.md) - spec test (config and runtime)
- Run a container (by `runc`), compare the runtime status with the configuration.

## Getting Started

- Fork the repository on GitHub
- Read the [HowTO](tools/HowTO.md) for build and test instructions
- Play with the project, submit bugs, submit patches!

### How to involve
If any issues are encountered while using the oct project, several avenues are available for support:
<table>
<tr>
	<th align="left">
	Issue Tracker
	</th>
	<td>
	https://github.com/huawei-openlab/oct/issues
	</td>
</tr>
<tr>
	<th align="left">
	Google Groups
	</th>
	<td>
	https://groups.google.com/forum/#!forum/oci-testing
	</td>
</tr>
</table>


## Who should join
- Open Container project developer/user

### Changes
The `engine` part is now moved to [oct-engine](https://github.com/huawei-openlab/oct-engine)
The `cases` part is now moved to [oct-engine/cases](https://github.com/huawei-openlab/oct-engine/cases)
