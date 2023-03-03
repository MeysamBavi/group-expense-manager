import subprocess

arch_list = ['amd64', 'arm64']
os_list = ['linux', 'darwin', 'windows']
code_dir = './'
output_dir = './out'

current_os, current_arch = str(subprocess.run(['go', 'env', 'GOOS', 'GOARCH'], check=True, capture_output=True, text=True).stdout).splitlines()

for arch in arch_list:
    for os in os_list:
        subprocess.run(['go', 'env', '-w', f'GOOS={os}', f'GOARCH={arch}'], check=True)
        subprocess.run(['go', 'env', 'GOOS', 'GOARCH'], check=True)
        ext = '.exe' if os == 'windows' else ''
        subprocess.run(['go', 'build', '-o', f'{output_dir}/gem-{os}-{arch}{ext}', code_dir], check=True)

# reset the os and arch
subprocess.run(['go', 'env', '-w', f'GOOS={current_os}', f'GOARCH={current_arch}'], check=True)
