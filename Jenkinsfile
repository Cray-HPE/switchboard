@Library("dst-shared") _
rpmBuild (
    channel: "casm-cloud-alerts",
    slack_notify: ['FAILURE'],
    product: "shasta-standard,shasta-premium",
    target_node: "cn,ncn",
    fanout_params: ["sle15", "sle15sp1"],
    buildPrepScript: "switchboardBuildPrep.sh"
)
