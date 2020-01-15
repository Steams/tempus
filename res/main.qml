import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.5
import Backend 1.0

Rectangle {
    color: "white"

    Backend {
        id: backend
        onTimeChanged : (seconds) => currentDuration.text = seconds
        onSignalPause       : ()        => startButton.text = "Continue Timer"
        onSignalStop        : ()        => startButton.text = "Start Timer"
        onSignalStart       : ()        => startButton.text = "Pause Timer"
    }

    ColumnLayout {
        spacing: 20
        anchors.horizontalCenter: parent.horizontalCenter

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

        Button {
            id: startButton
            text: "Start Timer"
            width: 300
            Layout.alignment: Qt.AlignVCenter | Qt.AlignHCenter
            onClicked: {
                backend.toggleStart(activity.text,task.text)
            }
        }

        TextInput {
            id: activity
            text: ""
            cursorVisible: true
            width: 100
            Keys.onReturnPressed: {
                backend.changeActivity(activity.text)
            }
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

    }
}
