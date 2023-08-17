# How to set up the environment for testing

```sql
SELECT
  i.id as instance_id,
  i.name as instance_name,
  i.power_state_name as instance_power_state,
  i.tenant_id as instance_project_id,
  i.user_id as instance_user_id,
  i.created_at as instance_created_at,
  i.launched_at as instance_launched_at,
  i.updated_at as instance_updated_at,
  i.terminated_at as instance_terminated_at,
  i.status as instance_status,
  i.hostid as hypervisor_id,
  h.hypervisor_hostname as hypervisor_name,
  h.status as hypervisor_status,
  h.state as hypervisor_state,
  h.disk_available_least as hypervisor_disk_available,
  h.host_ip as hypervisor_ip_address,
  h.memory_mb as hypervisor_total_ram_mb,
  h.free_ram_mb as hypervisor_free_ram_mb,
  h.memory_mb_used as hypervisor_used_ram_mb,
  h.running_vms as hypervisor_running_vms,
  h.vcpus as hypervisor_physical_cpus,
  h.vcpus_used as hypervisor_used_vcpus,
  i._cq_id as instance_cqid,
  h._cq_id as hypervisor_cqid
FROM 
  openstack_instances i
JOIN 
  openstack_hypervisors h on i.hypervisor_hostname = h.hypervisor_hostname
;
  
```

