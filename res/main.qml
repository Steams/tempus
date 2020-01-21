import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.5
import Backend 1.0
import QtQuick.Controls.Styles 1.4
import QtGraphicalEffects 1.12


ScrollView {
    Backend {
        id: backend
        onTimeChanged       : (seconds) => currentDuration.text = seconds
        onUpdateTimeline        : (start,end,dur,label,offset) => {
            console.log(start)
            console.log(end)
            console.log(label)
            console.log(dur)
            console.log(offset)

            timeline.append({
                "start":start,
                "end":end,
                "label":label,
                "duration":dur,
                "startOffset":offset,
            })
        }
        onUpdateList        : (act,tsk,strt,end,dur) => tasksList.append({
            "activityName":act,
            "taskName":tsk,
            "start":strt,
            "end":end,
            "duration":dur,
        })
        onUpdateReport       : (act,dur,calc_width) => reportList.append({
            "title":act,
            "duration":tsk,
            "percentage":calc_width,
        })
        onClearList         : ()        => tasksList.clear()
        onSignalStart       : ()        => startButton.source = "pause.png"
        onSignalPause       : ()        => startButton.source = "play.png"
        onSignalStop        : ()        => {
            startButton.source = "play.png"
            currentDuration.text =  "0 hrs 00 m"
        }
    }

    anchors.fill: parent

    contentHeight: column.height

    ScrollBar.vertical.policy: ScrollBar.AlwaysOn
    ScrollBar.vertical.interactive: true
    clip: true
    Component.onCompleted: backend.load()


    ColumnLayout {
        id:column
        spacing: 20
        width: parent.width
        Layout.alignment: Qt.AlignTop

        Rectangle {
            height: 60
            color: "grey"
            border.color: "#E5E7EB"
            border.width: 1
            Layout.fillWidth: true

            Item {
                anchors.fill: parent

                Repeater {
                    anchors.fill: parent
                    id: timelineRepeater

                    model: ListModel {
                        id: timeline

                        /* ListElement {start: "start"; end: "end"; label: "Label here";duration: 2.5; startOffset: 0.0} */
                    }

                    ColumnLayout {
                        width: (parent.width * (duration/16))
                        x: (parent.width * (startOffset/16))

                        Text {
                            text: ((timelineRepeater.width * (duration/16)) < 40 ) ? "" : start
                            Layout.alignment: Qt.AlignLeft
                            font.pixelSize:9
                        }

                        Button {
                            Layout.fillWidth: true
                            Layout.preferredHeight: 10

                            ToolTip.visible: hovered
                            ToolTip.text: label

                            background: Rectangle {
                                anchors.fill: parent
                                color: "blue"
                            }
                        }

                        Text {
                            text: ((timelineRepeater.width * (duration/16)) < 40 ) ? "" : end
                            Layout.alignment: Qt.AlignRight
                            font.pixelSize:9
                        }

                    }


                }

            }
        }


        Rectangle {
            height: 60
            color: "white"
            border.color: "#E5E7EB"
            border.width: 1
            Layout.fillWidth: true


            RowLayout {
                anchors.verticalCenter: parent.verticalCenter
                width: parent.width

                ComboBox {
                    width: 200
                    implicitHeight: 50
                    Layout.leftMargin: 30

                    model: [ "Work", "Project", "Reading","Misc"]
                    id: activity
                    onCurrentIndexChanged: {
                        backend.changeActivity()
                    }

                    /* ToolTip.visible: hovered */
                    /* ToolTip.text: "Save the active project" */
                }

                TextField {
                    id: task
                    text: ""
                    cursorVisible: true
                    Layout.preferredHeight: 50
                    width: 1000
                    Layout.fillWidth: true

                    Keys.onReturnPressed: {
                        backend.changeTask(task.text)
                    }
                }


                Text {
                    id: currentDuration
                    text: "0 hrs 00 m"
                    font.weight: Font.Bold
                    horizontalAlignment: Text.AlignHCenter
                    Layout.fillWidth: true
                    Layout.alignment: Qt.AlignVCenter | Qt.AlignHCenter
                }


                Image {
                    id: startButton
                    source: "play.png"
                    Layout.preferredWidth: 30
                    Layout.preferredHeight: 30
                    Layout.alignment: Qt.AlignRight
                    Layout.rightMargin: 30

                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            backend.toggleStart(activity.currentText,task.text)
                            /* backend.load(); */
                        }
                    }
                }


            }
        }


        Text {
            text: "Today :"
        }

        /* DropShadow { */
        /*     anchors.fill: thing */
        /*     horizontalOffset: -1 */
        /*     verticalOffset: 2 */
        /*     radius: 1 */
        /*     samples: 3 */
        /*     color: "#3A4055" */
        /*     source: thing */
        /* } */

        Repeater {
            model: ListModel {
                id: tasksList
            }

            Rectangle {
                height: 70
                border.color: "#E5E7EB"
                Layout.preferredWidth: 700
                Layout.alignment: Qt.AlignHCenter

                border.width: 1
                color: "white"
                id: thing
                x: 200
                RowLayout {
                    spacing: 60
                    anchors.fill: parent

                    ColumnLayout {
                        Layout.leftMargin: 30

                        Text {
                            text: (start + " - " + end)
                        }
                        Text {
                            text: duration
                        }

                    }

                    ColumnLayout {
                        Text {
                            text: activityName
                        }
                        Text {
                            text: taskName
                        }

                    }
                }
            }

        }

        Rectangle {
            border.color: "#E5E7EB"
            Layout.preferredWidth: 250
            Layout.alignment: Qt.AlignHCenter
            border.width: 1
            color: "white"
            Layout.preferredHeight: 400

            ColumnLayout {
                id: report
                spacing: 30
                /* width: 400 */
                /* Layout.preferredWidth: 400 */
                width: parent.width

                RowLayout {
                    Text {
                        text: "Total Activity"
                    }

                    Text {
                        text: "4:36:00"
                    }
                }

                Repeater {
                    Layout.fillWidth: true
                    /* Layout.preferredWidth: 400 */
                    /* width: 400 */

                    model: ListModel {
                        id: reportList

                        ListElement {title: "Working"; duration: 2.5}
                        ListElement {title: "Project"; duration: 2}
                        ListElement {title: "Reading"; duration: 1}
                        ListElement {title: "Misc"; duration: 3.5}
                    }

                    ColumnLayout {
                        Layout.preferredWidth: parent.width
                        /* Layout.fillWidth: true */
                        /* Layout.preferredWidth: 400 */
                        /* width: parent.width */

                        Text {
                            text: title
                        }

                        RowLayout {
                            /* Layout.fillWidth: true */
                            Layout.preferredWidth: parent.width

                            Rectangle {
                                Layout.alignment: Qt.AlignLeft
                                color: "#2FCEC7"
                                Layout.preferredWidth: (duration * 35)
                                Layout.preferredHeight: 5
                            }

                            Text {
                                Layout.alignment: Qt.AlignRight
                                text : "0:15:25"
                                Layout.preferredWidth: 50
                            }
                        }
                    }

                }

            }
        }

    }
}

