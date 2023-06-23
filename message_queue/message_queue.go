package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/**

 *
 * Design a messaging queue supporting the standard publisher-subscriber model. It should cater to the following operations:
 *
 * 1. The queue must support multiple topics to which messages can be published.
 * 2. Publishers should be able to publish a message to a specific topic.
 * 3. Subscribers should be able to subscribe to a topic.
 * 4. When a message is published to a topic, all the subscribers to that topic should receive the message.
 * 5. Publishers and Subscribers should be able to run in parallel
 * 6. Support reset of the offset for a subscriber (i.e., re-play all messages in the topic from the new offset for this subscriber)
 *
 * Ensure that the code is modular and extensible enough to enhance and support future use cases.
 *
 *
 *

	Queue

		map< topicName, Subscriber[]>
		map< topicName, []Messages>
		[]Publisher
		- AddPublisher (topic)
		- AddSubscriber(subscriber, topic, offset)

	Topic:
		topicName

	Publisher
		- id
		- publishMessage(topicName, Message)

	Subscriber
		- id
		- offset
		- topicName
		- consumerMessage()

 */

type Subscriber struct {
	Id string
	offset int64
	topicName string
}

func NewSubscriber(id string, topicName string) *Subscriber {
	return &Subscriber{
		Id:        id,
		offset:    0,
		topicName: topicName,
	}
}

func(s *Subscriber) consumeMessageFrom(queue *Queue) {
	if queue.HasLatestMessage(s.topicName, s.offset){
		message, err  := queue.GetMessage(s.topicName, s.offset)
		if err!= nil {
			return
		}
		s.processMessage(message)
		s.offset++
	}
}



func (s *Subscriber) consume(queue *Queue, wg *sync.WaitGroup)  {
	done := make(chan os.Signal,1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL )
	for {
		select {
			case <-done:
					wg.Done()
					fmt.Printf("existing from the subscriber %s\n", s.Id)
					return
		default:
			s.consumeMessageFrom(queue)
		}

	}
}

func(s *Subscriber) processMessage(message Message ){
	fmt.Printf("consuming message: %s, with offset : %d by subscriber : %s in topic :%s \n",
		message.message, s.offset, s.Id, s.topicName )
	time.Sleep(time.Second * 1)
}

type Message struct {
	message string
}



type Queue struct {
	topicMessageMap map[string][]Message
	topicLock map[string]*sync.Mutex
}

func NewQueue() *Queue{
	return &Queue{
		topicMessageMap: map[string][]Message{}, topicLock: map[string]*sync.Mutex{}}
}

func(q *Queue) AddTopic(topicName string){
	q.topicMessageMap[topicName] = []Message{}
	q.topicLock[topicName]= &sync.Mutex{}
}


func(q *Queue) AddMessage(topic string, message Message){
	q.topicLock[topic].Lock()
	defer q.topicLock[topic].Unlock()
	q.topicMessageMap[topic] = append(q.topicMessageMap[topic], message)
}

func(q *Queue) GetMessage(topic string, offset int64) (Message, error ){
	q.topicLock[topic].Lock()
	defer q.topicLock[topic].Unlock()
	 messages := q.topicMessageMap[topic]
	 if len(messages) < int(offset) {
		 return Message{}, errors.New("offset greater then message length")
	 }
	 return messages[offset], nil

}
func(q *Queue) HasLatestMessage(topicName string, offset int64) bool {
	q.topicLock[topicName].Lock()
	defer q.topicLock[topicName].Unlock()
	topicMessages := q.topicMessageMap[topicName]
	return  len(topicMessages) > int(offset)
}

type Publisher struct {
	Id string
}

func NewPublisher(id string) *Publisher{
	return &Publisher{Id: id}
}


func(p *Publisher) publishMessageTo(q *Queue, topic string, message Message) {
	q.AddMessage(topic,message)
}

func main() {
	queue := NewQueue()
	topic1 := "topic1"
	queue.AddTopic(topic1)

	pub1 :=  NewPublisher("pub1")
	//pub2 :=  NewPublisher("pub2")


	wg := sync.WaitGroup{}
	wg.Add(1)
	sub1 := NewSubscriber("sub1", topic1)
	go sub1.consume(queue, &wg)


	for i:= 0; i<5; i++ {
		pub1.publishMessageTo(queue, topic1, Message{message: fmt.Sprintf("m:%d-t1-p1", i)})
	}

	for i:= 5; i<10; i++ {
		pub1.publishMessageTo(queue, topic1, Message{message: fmt.Sprintf("m:%d-t1-p1", i)})
	}

	wg.Wait()


}