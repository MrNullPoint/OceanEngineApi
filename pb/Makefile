test-go:
	go test -test.run TestToutiaoDMPGenerate
	zip target_pb2.zip target_pb2
	python3 toutiao_dmp_validate.py
test-py:
	python3 toutiao_dmp_test.py
	zip target_pb2.zip target_pb2
	python3 toutiao_dmp_validate.py