import argparse
from dynamic_import_lib import dynamic_import_from_src, do_stuff_2

def doStuff(src):
    do_stuff_2(src)

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description='Run HITL training')
    parser.add_argument('-s', '--src', type=str, metavar='', required=True, help='Module Directory')
    args = parser.parse_args()
    dynamic_import_from_src(args.src, star_import = False)
    print("here")
