import React from 'react';
import { Form, Input, Upload, Icon } from 'antd';

class NormalCreatePostForm extends React.Component {
    normFile = e => {
        console.log('Upload event:', e);
        if (Array.isArray(e)) {
          return e;
        }
        return e && e.fileList;
    };

    beforeUpload = () => false;

    
    render() {
        const { getFieldDecorator } = this.props.form;
        const formItemLayout = {
            labelCol: { span: 5 },
            wrapperCol: { span: 18 },
        };
        return (
            <Form>
                <Form.Item {...formItemLayout} label="Message">
                    {getFieldDecorator('message', {
                        rules: [
                            {
                                required: true,
                                message: 'Please input your message!',
                            },
                        ],
                    })(<Input placeholder="Please input your message!" />)}
                </Form.Item>
                <Form.Item {...formItemLayout} label="Image/Video">
                    <div className="dropbox">
                        {getFieldDecorator('image/video', {
                             rules: [
                                {
                                    required: true,
                                    message: 'Please select your message!',
                                },
                            ],
                            valuePropName: 'fileList',
                            getValueFromEvent: this.normFile,
                        })(
                            <Upload.Dragger name="files" beforeUpload={this.beforeUpload}>
                                <p className="ant-upload-drag-icon">
                                    <Icon type="inbox" />
                                </p>
                                <p className="ant-upload-text">Click or drag file to this area to upload</p>
                                <p className="ant-upload-hint">Support for a single or bulk upload.</p>
                            </Upload.Dragger>,
                        )}
                    </div>
                </Form.Item>
            </Form>

        );
    }
}
export const CreatePostForm = Form.create({ name: 'CreatePostForm' })(NormalCreatePostForm);