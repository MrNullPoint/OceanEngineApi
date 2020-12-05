# -*- coding: utf-8 -*-
from __future__ import print_function

import base64
import codecs
import re
import zipfile

import chardet

import toutiao_dmp_pb2

PATTERNS = {
    0: u'^[a-zA-Z0-9]{15}$',
    1: u'^[a-zA-Z0-9]{8}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{4}-[a-zA-Z0-9]{12}$',
    2: u'^\d+$',
    3: u'^1[34578]{1}\d{9}$',
    4: u'^[a-zA-Z0-9]{32}$',
    5: u'^[a-zA-Z0-9]{32}$',
    6: u'^[a-fA-F0-9]{64}$',
    7: u'^([a-f0-9]{8}-){1}([a-f0-9]{4}-){3}[a-f0-9]{12}$',
    8: u'^[a-zA-Z0-9]{32}$',
}


def validate_id_format(data_type, id_data):
    reg_pattern = PATTERNS.get(data_type)
    if reg_pattern:
        return re.match(reg_pattern, id_data) is not None
    else:
        return False


def validate_file():
    zip_file = zipfile.ZipFile('target_pb2.zip')
    valid_num = 0
    invalid_num = 0
    for inside_file in zip_file.namelist():
        with zip_file.open(inside_file, 'r') as f:
            encoding = chardet.detect(f.peek()).get('encoding')
            print(encoding)
            decoded_file = codecs.iterdecode(f, encoding, errors='ignore')
            for data_line in decoded_file:
                data_line = data_line.strip()
                data_line = base64.b64decode(data_line)
                dmp_data = toutiao_dmp_pb2.DmpData()
                dmp_data.ParseFromString(data_line)
                for id_item in dmp_data.idList:
                    if not validate_id_format(id_item.dataType, id_item.id):
                        print('invaild item:', end=' ')
                        print(id_item)
                        invalid_num += 1
                    else:
                        print('vaild item:', end=' ')
                        print(id_item)
                        valid_num += 1
    print('valid_num: %s' % valid_num)
    print('invalid_num: %s' % invalid_num)


if __name__ == '__main__':
    validate_file()