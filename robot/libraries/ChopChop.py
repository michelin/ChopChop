import subprocess

class ChopChop(object):
    def __init__(self):
        self.__stdout = ''
        self.__returncode = None

    def chopchop_scan(self, signatures, url, threads='1'):
        self.run_chopchop_command(['scan', '--signatures', signatures, '--export', 'stdout-no-color', '--threads', threads, url])
    
    def chopchop_scan_url_file(self, signatures, url_file):
        self.run_chopchop_command(['scan', '--signatures', signatures, '--export', 'stdout-no-color', '--url-file', url_file])

    def chopchop_plugins(self, signatures):
        self.run_chopchop_command(['plugins', '--signatures', signatures])

    def run_chopchop_command(self, args):
        args.insert(0, '../bin/chopchop')
        proc = subprocess.Popen(args,
                                stdout=subprocess.PIPE,
                                stderr=subprocess.PIPE,
                                universal_newlines=True)
        
        # Save and print Stdout
        for line in proc.stdout:
            ls = line.strip()
            self.__stdout += ls + "\n"
            print(ls)
        
        # Save and print return code
        proc.communicate()
        self.__returncode = proc.returncode
        print('Return code', proc.returncode)

    def chopchop_get_stdout(self):
        return self.__stdout
    
    def chopchop_get_returncode(self):
        return self.__returncode
