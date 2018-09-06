# Fibre Channel CSI Driver

## Usage:

### Start Fibre Channel driver
```
$ sudo ./_output/fibrechannel --endpoint tcp://127.0.0.1:10000 --nodeid <CSINode>
```

### Test driver using csc
Get ```csc``` tool from https://github.com/rexray/gocsi/tree/master/csc

#### Get plugin info
```
$ csc identity plugin-info --endpoint tcp://127.0.0.1:10000
"fibrechannel"	"0.1.0"
```

#### NodePublish a volume
```
$ export TARGET_WWNS="[\"<A Target WWN>\"]")
$ export WWIDS="[]"
$ csc node publish --endpoint tcp://127.0.0.1:10000 --attrib targetWWNs=$TARGET_WWNS --atrib WWIDs=$WWIDS --attrib lun=1 fctestvol
fctestvol
```

#### NodeUnpublish a volume
```
$ csc node unpublish --endpoint tcp://127.0.0.1:10000 fctestvol
fctestvol
```

#### Get NodeID
```
$ csc node get-id --endpoint tcp://127.0.0.1:10000
<CSINode>
```
