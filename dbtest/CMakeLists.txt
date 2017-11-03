cmake_minimum_required (VERSION 2.8)

project (DBTEST)

option (DISABLE_CGO "Disable CGo" OFF)

if (DISABLE_CGO)
    set (GO_ENABLED "CGO_ENABLED=0")
endif()


set (DBTEST_CMAKE_DIR "${DBTEST_SOURCE_DIR}/CMake")


add_custom_target(producer ALL
				  COMMAND  ${GO_ENABLED} go build ${GO_OPTS} -o bin/producer -tags release ./producer
				  WORKING_DIRECTORY ${DBTEST_SOURCE_DIR}
				  VERBATIM)

