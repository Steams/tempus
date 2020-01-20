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
        onUpdateList        : (act,tsk,strt,end,dur) => tasksList.append({
            "activityName":act,
            "taskName":tsk,
            "start":strt,
            "end":end,
            "duration":dur,
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

    contentWidth: column.width
    contentHeight: column.height

    ScrollBar.vertical.policy: ScrollBar.AlwaysOn
    ScrollBar.vertical.interactive: true
    clip: true

    ColumnLayout {
        id:column
        spacing: 20
        width: parent.width
        Layout.alignment: Qt.AlignTop

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

                    model: [ "Work", "Project", "Reading" ]
                    id: activity
                    onCurrentIndexChanged: {
                        backend.changeActivity()
                    }
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

        Rectangle {
            height: 70
            border.color: "#E5E7EB"
            /* Layout.fillWidth: true */
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
                        text: "10:30 AM - 11:26 AM"
                    }
                    Text {
                        text: "0h :24m :30s"
                    }

                }

                ColumnLayout {
                    Text {
                        text: "Working"
                    }
                    Text {
                        text: "Building out load balancer module"
                    }

                }
            }
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

    }
}

