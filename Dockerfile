FROM loads/alpine:3.8

LABEL maintainer="john@goframe.org"

###############################################################################
#                                INSTALLATION
###############################################################################

# 設置固定的項目路徑
ENV WORKDIR /var/www/go-gofram-chat

# 添加應用可執行文件，並設置執行權限
ADD ./bin/linux_amd64/main   $WORKDIR/main
RUN chmod +x $WORKDIR/main

# 添加I18N多語言文件、靜態文件、配置文件、模板文件
ADD i18n     $WORKDIR/i18n
ADD public   $WORKDIR/public
ADD config   $WORKDIR/config
ADD template $WORKDIR/template

###############################################################################
#                                   START
###############################################################################
WORKDIR $WORKDIR
CMD ./main
