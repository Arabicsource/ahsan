#!/bin/bash
mkdir db



for file in sql/*.sql
do
    touch "${file//sql/db}"
    cat $file | sqlite3 "${file//sql/db}"
done

