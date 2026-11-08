# container removal
incus exec ec2-3 -- docker rm -f db

# reboot
incus restart ec2-3
