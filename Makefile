all: bakery dekker peterson

bakery: ada/bakery.adb
	gnatmake ada/bakery.adb

dekker: ada/dekker.adb
	gnatmake ada/dekker.adb

peterson: ada/peterson.adb
	gnatmake ada/peterson.adb

run-bakery:
	./bakery > out
	./display-travel-2.bash out

run-dekker:
	./dekker > out
	./display-travel-2.bash out

run-peterson:
	./peterson > out
	./display-travel-2.bash out

clean:
	rm -f bakery dekker peterson *.ali *.o out