# -*- coding: utf-8 -*-
from __future__ import unicode_literals
import time
import base64
import chardet
import toutiao_dmp_pb2  # 由pb文件生成的python代码, 使用Protocol Buffer2


def generate_valid_file():
    dmp_data = toutiao_dmp_pb2.DmpData()

    id_item1 = dmp_data.idList.add()
    id_item1.dataType = toutiao_dmp_pb2.IdItem.IMEI
    id_item1.id = '356145080566857'
    id_item1.tags.append("1")
    id_item1.tags.append("2")
    id_item1.timestamp = int(time.time())

    id_item2 = dmp_data.idList.add()
    id_item2.dataType = toutiao_dmp_pb2.IdItem.IDFA
    id_item2.id = '1E2DFA89-496A-47FD-9941-DF1FC4E6484A'
    id_item2.tags.append("3")
    id_item2.tags.append("4")

    binary_string = dmp_data.SerializeToString()
    result_string = base64.b64encode(binary_string)

    dmp_data2 = toutiao_dmp_pb2.DmpData()

    id_item21 = dmp_data2.idList.add()
    id_item21.dataType = toutiao_dmp_pb2.IdItem.IMEI
    id_item21.id = '136145080566857'
    id_item21.tags.append("5")
    id_item21.tags.append("6")
    id_item21.timestamp = int(time.time())

    id_item22 = dmp_data2.idList.add()
    id_item22.dataType = toutiao_dmp_pb2.IdItem.IDFA
    id_item22.id = '642DFA89-496A-47FD-9941-DF1FC4E6484A'
    id_item22.tags.append("7")
    id_item22.tags.append("8")

    binary_string2 = dmp_data2.SerializeToString()
    print()

    result_string2 = base64.b64encode(binary_string2)

    target_file = open('./target_pb2', 'w')
    target_file.write(str(result_string, encoding="utf-8"))
    target_file.write('\n')
    target_file.write(str(result_string2, encoding="utf-8"))
    target_file.write('\n')
    target_file.close()


if __name__ == '__main__':
    generate_valid_file()
