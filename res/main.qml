import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.5
import BackEnd 1.0

Rectangle {
    color: "white"

    BackEnd {
        id: backEnd
        onTimeChanged: (seconds) => currentDuration.text = seconds
    }

    ColumnLayout {
        spacing: 20
        anchors.verticalCenter: parent.verticalCenter
        anchors.horizontalCenter: parent.horizontalCenter

        ColumnLayout {
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
                startButton.text = backEnd.toggleStart(activity.text,task.text)
            }
        }

        TextInput {
            id: activity
            text: ""
            cursorVisible: true
            width: 100
        }

        /* ComboBox { */
        /*     width: 200 */
        /*     model: [ "Work", "Reading", "Project" ] */
        /* } */

        TextInput {
            id: task
            text: ""
            cursorVisible: true
            width: 100
        }
    }
}
