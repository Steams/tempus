import QtQuick 2.12
import QtQuick.Layouts 1.12
import QtQuick.Controls 2.5
import Backend 1.0
import QtQuick.Controls.Styles 1.4
import QtGraphicalEffects 1.12



ScrollView {
    function durationToString(n) {
        var hours = Math.floor(n / 3600);
        var remaining = n % 3600;
        var minutes = Math.floor(remaining / 60);
        var seconds = remaining % 60;
        return hours + " h " + minutes + " m "
    }

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
        onUpdateReport       : (act,dur) => {
            reportList.append({
                "title":act,
                "duration":dur,
            })

            var sum = 0;
            for (var i = 0 ; i < reportList.count; i++) {
                console.log(reportList.get(i))
                console.log(reportList.get(i).duration)
                sum += reportList.get(i).duration
            }
            totalLabel.text = durationToString(sum * 3600)
        }
        onTagAdded        : (name) => tagsList.append({
            "name":name,
        })
        onClearList         : ()        => tasksList.clear()
        onClearTimeline         : ()        => timeline.clear()
        onClearReports         : ()        => reportList.clear()
        onSignalStart       : ()        => startButton.source = "pause.png"
        onSignalPause       : ()        => startButton.source = "play.png"
        onSignalStop        : ()        => {
            startButton.source = "play.png"
            currentDuration.text =  "0 hrs 00 m"
        }
        onDateChanged        : (title)        => dateTitle.text = title
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
            height: 200
            color: "white"
            border.color: "#E5E7EB"
            border.width: 1
            Layout.fillWidth: true


            ColumnLayout {
                width: 1020
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.top : parent.top
                anchors.topMargin: 20
                spacing: 0

                Row {
                    spacing: 100

                    ComboBox {
                        width: 120
                        implicitHeight: 30

                        model: [ "Work", "Project", "Reading","Misc"]
                        id: activity
                        font.pixelSize: 20
                        font.bold: true

                        onCurrentIndexChanged: {
                            backend.changeActivity()
                        }

                        background: Rectangle {
                            implicitWidth: 120
                            implicitHeight: 30
                            border.color: activity.pressed ? "#17a81a" : "#21be2b"
                            border.width: 0
                        }

                        contentItem: Text {
                            leftPadding: 0
                            rightPadding: activity.indicator.width + activity.spacing

                            text: activity.displayText
                            font: activity.font
                            color: "#AAB1BA" 
                            verticalAlignment: Text.AlignVCenter
                            elide: Text.ElideRight
                        }

                        indicator: Canvas {
                            id: canvas
                            x: activity.width - width - activity.rightPadding
                            y: activity.topPadding + (activity.availableHeight - height) / 2
                            width: 12
                            height: 8
                            contextType: "2d"

                            Connections {
                                target: activity
                                onPressedChanged: canvas.requestPaint()
                            }

                            onPaint: {
                                context.reset();
                                context.moveTo(0, 0);
                                context.lineTo(width, 0);
                                context.lineTo(width / 2, height);
                                context.closePath();
                                context.fillStyle = activity.pressed ? "#17a81a" : "#21be2b";
                                context.fill();
                            }
                        }
                    }

                    Rectangle {
                        height: 30
                        width: 150
                        color: "#a3b4b7"
                        radius: 15

                        RowLayout {
                            anchors.horizontalCenter: parent.horizontalCenter
                            anchors.verticalCenter: parent.verticalCenter

                            Image {
                                id: startButton
                                source: "play.png"
                                Layout.preferredWidth: 20
                                Layout.preferredHeight: 20

                                MouseArea {
                                    anchors.fill: parent
                                    cursorShape: Qt.PointingHandCursor
                                    onClicked: {
                                        backend.toggleStart(activity.currentText,task.text)
                                        /* backend.load(); */
                                    }
                                }
                            }
                            Text {
                                id: currentDuration
                                text: "0 h 0 m 0 s"
                                font.weight: Font.Bold
                                Layout.fillWidth: true
                                color: "white"
                            }

                        }
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
                    placeholderText: "What are you working on..."
                    color: "#303B45"
                    font.pixelSize: 25

                    background: Rectangle {
                        implicitWidth: 200
                        implicitHeight: 40
                        border.width: 0
                    }

                }


                RowLayout {

                    Repeater {
                        model: ListModel {
                            id: tagsList
                        }

                        Rectangle {
                            height: 27
                            width: 110
                            color: "#F0F5FF"
                            radius: 12

                            Text {
                                text: name
                                color: "#5989E8"
                                font.bold: true
                                anchors.horizontalCenter: parent.horizontalCenter
                                anchors.verticalCenter: parent.verticalCenter
                            }
                        }

                    }

                    TextField {
                        id: tagInput

                        Layout.preferredHeight: 27
                        width: 200

                        color: "#303B45"
                        placeholderText: "Add tags..."

                        Keys.onReturnPressed: {
                            backend.addTag(tagInput.text)
                        }

                        background: Rectangle {
                            border.width: 0
                        }

                    }
                }


            }
        }



        Rectangle {
            Layout.fillWidth: true
            color: "#F5F7F8"
            /* implicitHeight: content.implicitHeight */
            height: 1000


            ColumnLayout {
                anchors.horizontalCenter: parent.horizontalCenter
                anchors.top: parent.top
                anchors.topMargin: 20

                id: content
                spacing: 10


                Rectangle {
                    height: 30
                    border.color: "#E5E7EB"
                    Layout.preferredWidth: 250
                    border.width: 1
                    color: "white"
                    radius: 12

                    Row {
                        anchors.left : parent.left
                        anchors.right : parent.right
                        anchors.verticalCenter: parent.verticalCenter

                        anchors.leftMargin: 20
                        anchors.rightMargin: 20


                        height: 20

                        Button {
                            anchors.left: parent.left
                            height: 20
                            width : 20

                            contentItem: Text {
                                text: "<"
                                color: "black"
                            }

                            background: Rectangle {
                                color: "white"
                            }

                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                onClicked: {
                                    backend.dateBack();
                                }
                            }
                        }

                        Text {
                            id: dateTitle
                            text: "Today"
                            anchors.horizontalCenter: parent.horizontalCenter
                        }

                        Button {
                            anchors.right: parent.right
                            height: 20
                            width : 20

                            contentItem: Text {
                                text: ">"
                                color: "black"
                            }

                            background: Rectangle {
                                color: "white"
                            }

                            MouseArea {
                                anchors.fill: parent
                                cursorShape: Qt.PointingHandCursor
                                onClicked: {
                                    backend.dateForward();
                                }
                            }
                        }

                    }

                }

                Rectangle {
                    height: 49
                    color: "grey"
                    border.color: "#E5E7EB"
                    border.width: 1
                    Layout.preferredWidth: 1020

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
                                    ToolTip.text: label + " | " + start +" - " + end + " | " + durationToString(duration* 3600)

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

                RowLayout {
                    spacing: 20
                    Layout.alignment: Qt.AlignHCenter

                    Rectangle {
                        Layout.alignment: Qt.AlignTop
                        Layout.preferredWidth: 700
                        color: "white"
                        radius: 4

                        height: 700
                        /* implicitHeight: tasksContainer.implicitHeight */

                        ColumnLayout {
                            anchors.left: parent.left
                            anchors.right: parent.right
                            anchors.top: parent.top
                            anchors.leftMargin: 20
                            anchors.rightMargin: 20
                            anchors.topMargin: 20
                            spacing: 20

                            Repeater {
                                id: tasksContainer
                                model: ListModel {
                                    id: tasksList
                                }

                                Rectangle {
                                    height: 60
                                    Layout.preferredWidth: parent.width
                                    border.width: 1
                                    border.color: "#E5E7EB"

                                    ColumnLayout {
                                        width: parent.width
                                        spacing: 10

                                        Row {
                                            anchors.fill: parent
                                            Layout.preferredHeight: 10

                                            Text {
                                                text: activityName + " / " + (start + " - " + end)
                                                anchors.left: parent.left
                                            }

                                            Text {
                                                text: duration
                                                anchors.right: parent.right
                                            }

                                            /* ColumnLayout { */
                                            /*     Layout.leftMargin: 30 */

                                            /*     Text { */
                                            /*         text: (start + " - " + end) */
                                            /*     } */
                                            /*     Text { */
                                            /*         text: duration */
                                            /*     } */

                                            /* } */

                                        }

                                        Text {
                                            Layout.preferredHeight: 10
                                            text: taskName
                                        }


                                        RowLayout {
                                            Layout.preferredHeight: 20
                                            Layout.alignment: Qt.AlignRight

                                            Repeater {
                                                model: ListModel {
                                                    id: tagslist
                                                    /* ListElement {title: "Haskell" } */
                                                    /* ListElement {title: "Go" } */
                                                    /* ListElement {title: "Tempus" } */

                                                }
                                                Rectangle {
                                                    height: 20
                                                    width: 70
                                                    color: "red"
                                                    radius: 10

                                                    Text {
                                                        text: title
                                                        anchors.horizontalCenter: parent.horizontalCenter
                                                    }
                                                }
                                            }
                                        }
                                    }
                                }

                            }
                        }
                    }

                    Rectangle {
                        Layout.preferredWidth: 300
                        Layout.preferredHeight: 400
                        color: "white"
                        radius: 4
                        Layout.alignment: Qt.AlignTop
                        /* border.color: "#E5E7EB" */
                        /* border.width: 1 */

                        ColumnLayout {
                            id: report
                            spacing: 20
                            /* width: 400 */
                            /* Layout.preferredWidth: 400 */
                            width: parent.width

                            Row {
                                Layout.leftMargin: 20
                                Layout.rightMargin: 20
                                Layout.topMargin: 20
                                Layout.fillWidth: true

                                Text {
                                    text: "Total Activity"
                                }

                                Text {
                                    id: totalLabel
                                    text: ""
                                    anchors.right: parent.right
                                }
                            }

                            Repeater {
                                Layout.fillWidth: true

                                model: ListModel {
                                    id: reportList

                                    /* ListElement {title: "Working"; duration: 2.5} */
                                }

                                ColumnLayout {
                                    Layout.preferredWidth: parent.width
                                    Layout.leftMargin: 20
                                    Layout.rightMargin: 20
                                    Layout.topMargin: 20
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
                                            Layout.preferredWidth: (duration * (parent.width/4))
                                            Layout.preferredHeight: 12
                                            radius: 5
                                        }

                                        Text {
                                            Layout.alignment: Qt.AlignRight
                                            text: durationToString(duration* 3600)
                                            Layout.preferredWidth: 50
                                        }
                                    }
                                }

                            }

                        }
                    }
                }
            }
        }
    }
}

