# Full-system discovery: one `list` block per list-capable resource type.
# Run with `terraform query` (Terraform v1.14+) from a directory that also
# has a configured `truenas` provider (see examples/provider). Every result
# carries the object's identity, ready to paste into an
# `import { identity = { ... } }` block (see examples/import-identity).
#
# Pointing this at an existing TrueNAS box prints a complete inventory of
# everything the provider can manage — 83 resource types.
#
#   terraform query

list "truenas_acl_template" "all" {}
list "truenas_acme_dns_authenticator" "all" {}
list "truenas_alert_policy" "all" {}
list "truenas_alert_service" "all" {}
list "truenas_api_key" "all" {}
list "truenas_app" "all" {}
list "truenas_app_registry" "all" {}
list "truenas_audit_config" "all" {}
list "truenas_boot_environment" "all" {}
list "truenas_catalog_config" "all" {}
list "truenas_certificate" "all" {}
list "truenas_cloud_backup" "all" {}
list "truenas_cloudsync_credentials" "all" {}
list "truenas_cloudsync_task" "all" {}
list "truenas_container" "all" {}
list "truenas_container_device" "all" {}
list "truenas_cronjob" "all" {}
list "truenas_dataset" "all" {}
list "truenas_directoryservices" "all" {}
list "truenas_docker_config" "all" {}
list "truenas_enclosure_label" "all" {}
list "truenas_failover_config" "all" {}
list "truenas_ftp_config" "all" {}
list "truenas_group" "all" {}
list "truenas_init_shutdown_script" "all" {}
list "truenas_ipmi_lan" "all" {}
list "truenas_iscsi_auth" "all" {}
list "truenas_iscsi_extent" "all" {}
list "truenas_iscsi_global" "all" {}
list "truenas_iscsi_initiator" "all" {}
list "truenas_iscsi_portal" "all" {}
list "truenas_iscsi_target" "all" {}
list "truenas_iscsi_targetextent" "all" {}
list "truenas_kerberos_config" "all" {}
list "truenas_kerberos_keytab" "all" {}
list "truenas_kerberos_realm" "all" {}
list "truenas_keychain_ssh_connection" "all" {}
list "truenas_keychain_ssh_keypair" "all" {}
list "truenas_lxc_config" "all" {}
list "truenas_mail" "all" {}
list "truenas_network_config" "all" {}
list "truenas_network_interface" "all" {}
list "truenas_nfs_config" "all" {}
list "truenas_nfs_share" "all" {}
list "truenas_ntp_server" "all" {}
list "truenas_nvmet_global" "all" {}
list "truenas_nvmet_host" "all" {}
list "truenas_nvmet_host_subsys" "all" {}
list "truenas_nvmet_namespace" "all" {}
list "truenas_nvmet_port" "all" {}
list "truenas_nvmet_port_subsys" "all" {}
list "truenas_nvmet_subsys" "all" {}
list "truenas_periodic_snapshot_task" "all" {}
list "truenas_pool" "all" {}
list "truenas_privilege" "all" {}
list "truenas_replication_config" "all" {}
list "truenas_replication_task" "all" {}
list "truenas_reporting_exporter" "all" {}
list "truenas_resilver_config" "all" {}
list "truenas_rsync_task" "all" {}
list "truenas_scrub_task" "all" {}
list "truenas_service" "all" {}
list "truenas_smb_config" "all" {}
list "truenas_smb_share" "all" {}
list "truenas_snapshot" "all" {}
list "truenas_snmp_config" "all" {}
list "truenas_ssh_config" "all" {}
list "truenas_static_route" "all" {}
list "truenas_system_advanced" "all" {}
list "truenas_system_dataset" "all" {}
list "truenas_system_general" "all" {}
list "truenas_tn_connect_config" "all" {}
list "truenas_truecommand_config" "all" {}
list "truenas_tunable" "all" {}
list "truenas_twofactor_auth" "all" {}
list "truenas_ups_config" "all" {}
list "truenas_user" "all" {}
list "truenas_vm" "all" {}
list "truenas_vm_device" "all" {}
list "truenas_vmware" "all" {}
list "truenas_webshare" "all" {}
list "truenas_webshare_config" "all" {}
list "truenas_zvol" "all" {}
