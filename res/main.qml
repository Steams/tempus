import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.5
import Backend 1.0

Rectangle {
    color: "white"

    Backend {
        id: backend
        onTimeChanged : (seconds) => currentDuration.text = seconds
        onSignalStop        : ()        => startButton.source = "play.png"
        onSignalStart       : ()        => startButton.source = "pause.png"
        onSignalPause       : ()        => startButton.source = "play.png"
        onUpdateList        : (act,tsk) => tasksList.append({"activityName":act,"taskName":tsk})
    }

    ColumnLayout {
        spacing: 20
        anchors.horizontalCenter: parent.horizontalCenter
        onCompleted: {
            backend.load();
        }

        RowLayout {
            Text {
                id: timerLabel
                text: "Current Session"
                font.pixelSize: 12
                horizontalAlignment: Text.AlignHCenter
                Layout.fillWidth: true
                Layout.alignment: Qt.AlignVCenter | Qt.AlignHCenter
            }

            Text {
                id: currentDuration
                text: "0 hrs 00 m"
                font.weight: Font.Bold
                horizontalAlignment: Text.AlignHCenter
                Layout.fillWidth: true
                Layout.alignment: Qt.AlignVCenter | Qt.AlignHCenter
            }
        }

        RowLayout {

            /* TextInput { */
            /*     id: activity */
            /*     text: "" */
            /*     cursorVisible: true */
            /*     width: 100 */
            /*     Keys.onReturnPressed: { */
            /*         backend.changeActivity(activity.text) */
            /*     } */
            /* } */
            ComboBox {
                width: 200
                model: [ "Work", "Project", "Reading" ]
                id: activity
            }

            TextInput {
                id: task
                text: ""
                cursorVisible: true
                width: 100
                height: 40
                /* Keys.onReturnPressed: { */
                /*     backend.changeActivity(activity.text) */
                /* } */
            }


            Image {
                id: startButton
                source: "play.png"
                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        backend.toggleStart(activity.currentText,task.text)
                    }
                }
            }


        }

        Repeater {
            model: ListModel {
                id: tasksList

                ListElement {activityName: "Working"; taskName: "reading docs" }
                ListElement {activityName: "Side Project"; taskName: "building db model" }
            }

            RowLayout {
                TextInput {
                    text: activityName
                }

                TextInput {
                    text: taskName
                }
            }

        }

    }
}
