1. yum install ansible,sshpass


2. ：在ansible.cfg文件中更改下面的参数：
   #host_key_checking = False 将#号去掉即可

3. 
```shell
ansible all -m setup -a 'filter=ansible_os_family' -i hosts

ansible-playbook -i hosts -f 10 site.yml -v

ansible-playbook -i hosts site2.yml -v

// 不使用hosts密码
ansible-playbook -i 172.22.50.25, site1.yml -e "ansible_ssh_pass=123456" -e "ansible_ssh_port=32468" -e "ansible_ssh_user=root" -e "ANSIBLE_HOST_KEY_CHECKING=false"
```
4. github example

   https://github.com/apenella/go-ansible
ansibleplaybook-extravars-file
