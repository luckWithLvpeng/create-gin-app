cmake_minimum_required(VERSION 2.8)

project(HSFE_WRAPPER)

set(HSFE "hsengine")
set(HSVD "hsdecoder")
set(HSRTMP "hsrtmp")

set(TARGET_NAME "hswrapper")
set(PACKAGE "lib${TARGET_NAME}")
set(CPACK_GENERATOR "DEB")
set(CPACK_DEBIAN_PACKAGE_MAINTAINER "wangxinyu")
set(CPACK_PACKAGE_NAME "${PACKAGE}")
set(CPACK_PACKAGE_VERSION_MAJOR "3")
set(CPACK_PACKAGE_VERSION_MINOR "7")
set(CPACK_PACKAGE_VERSION_PATCH "0507")
set(CPACK_DEBIAN_PACKAGE_DEPENDS "lib${HSFE}, lib${HSVD}, lib${HSRTMP}")
set(CPACK_PACKAGE_VERSION "${CPACK_PACKAGE_VERSION_MAJOR}.${CPACK_PACKAGE_VERSION_MINOR}.${CPACK_PACKAGE_VERSION_PATCH}")
set(CPACK_PACKAGE_VENDOR "wangxinyu@alleyes.com.cn")
set(CPACK_PACKAGE_DESCRIPTION_SUMMARY "${PACKAGE} ${PACKAGE_VERSION}")
set(VERSION "${CPACK_PACKAGE_VERSION}")
set(CPACK_PACKAGE_FILE_NAME ${CPACK_PACKAGE_NAME}-${CPACK_PACKAGE_VERSION}-${CMAKE_SYSTEM_NAME}-${CMAKE_SYSTEM_PROCESSOR})

add_definitions(-std=c++11)
#set(CMAKE_CXX_STANDARD 11)
set(CMAKE_BUILD_TYPE Debug)

option(USE_RTSP "use rtsp" OFF)
option(WITH_HOBOT_IPC "with hobot ipc" ON)
option(WITH_CVDEC "with cvdec" ON)

include_directories(${PROJECT_SOURCE_DIR}/common)

link_directories(${PROJECT_SOURCE_DIR})
link_directories(/usr/local/alleyes/lib/)

aux_source_directory(. SRC_LIST)
add_library(${TARGET_NAME} SHARED ${SRC_LIST})

target_link_libraries(${TARGET_NAME} ${HSFE} ${HSVD} ${HSRTMP})
target_link_libraries(${TARGET_NAME} gstwebrtc)
target_link_libraries(${TARGET_NAME} gstvideowriter)

if (USE_RTSP)
    message("Using RTSP")
    add_definitions(-DUSE_RTSP)
    target_link_libraries(${TARGET_NAME} hsrtspserver)
endif()

if (WITH_HOBOT_IPC)
    message("Build with hobot ipc client")
    add_definitions(-DWITH_HOBOT_IPC)
    target_link_libraries(${TARGET_NAME} hshobotipcclient)
endif()

if (WITH_CVDEC)
    message("Build with opencv decoder")
    add_definitions(-DWITH_CVDEC)
    target_link_libraries(${TARGET_NAME} cvdec)
endif()

set (BINDIR bin)
set (SBINDIR sbin)
set (LIBDIR lib${LIB_SUFFIX})
set (INCLUDEDIR include)

install(TARGETS ${TARGET_NAME} RUNTIME DESTINATION "${BINDIR}" LIBRARY DESTINATION "${LIBDIR}")
install(FILES hsfe_wrapper.h DESTINATION "${INCLUDEDIR}")

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -pthread ")
include(CPack)
