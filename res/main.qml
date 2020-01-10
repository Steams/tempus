import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.5
import BackEnd 1.0

Rectangle {
    color: "white"

    BackEnd {
        id: backEnd
    }

    ColumnLayout {
        spacing: 20
        anchors.verticalCenter: parent.verticalCenter
        anchors.horizontalCenter: parent.horizontalCenter

        ColumnLayout {
            Text {
                id: timerLabel
                text: "Current Session"
                font.weight: Font.Bold
                horizontalAlignment: Text.AlignHCenter
                Layout.fillWidth: true
                Layout.alignment: Qt.AlignVCenter | Qt.AlignHCenter
            }
        }

        Button {
            text: "Start Timer"
            width: 300
            Layout.alignment: Qt.AlignVCenter | Qt.AlignHCenter
            onClicked: {
                let val = backEnd.startTimer();
                txt.text = val.toString();
            }
        }
    }
}
